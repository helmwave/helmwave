package release

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	gotemplate "text/template"

	"github.com/google/shlex"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/invopop/jsonschema"
	cp "github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/postrender"
	"sigs.k8s.io/kustomize/api/konfig"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// PostRendererReference supports exec and kustomize modes for post-rendering.
// Only one of Exec or Kustomize should be set.
//
//nolint:lll
type PostRendererReference struct {
	Kustomize *KustomizeConfig `yaml:"kustomize,omitempty" json:"kustomize,omitempty" jsonschema:"description=Kustomize directory for post-rendering"`
	Exec      ExecConfig       `yaml:"exec,omitempty" json:"exec,omitempty" jsonschema:"description=Command to execute for post-rendering"`
}

// ExecConfig represents a command to execute.
// Supports both string (will be split using shlex) and array formats.
type ExecConfig []string

// KustomizeConfig specifies a kustomization directory for post-rendering.
//
//nolint:lll
type KustomizeConfig struct {
	Src            string `yaml:"src" json:"src" jsonschema:"required,description=Source directory path containing kustomization.yaml"`
	Dst            string `yaml:"dst" json:"dst" jsonschema:"readOnly,description=Rendered destination directory"`
	DelimiterLeft  string `yaml:"delimiter_left,omitempty" json:"delimiter_left,omitempty" jsonschema:"Set left delimiter for template engine,default={{"`
	DelimiterRight string `yaml:"delimiter_right,omitempty" json:"delimiter_right,omitempty" jsonschema:"Set right delimiter for template engine,default=}}"`
	Renderer       string `yaml:"renderer" json:"renderer" jsonschema:"description=How to render the files,enum=sprig,enum=gomplate,enum=copy,default=sprig"`
}

func (p PostRendererReference) JSONSchema() *jsonschema.Schema {
	// Use reflector for the object format, and manually add array format as OneOf
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
	}

	type postRendererReference PostRendererReference
	objectSchema := r.Reflect(postRendererReference(p))

	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "array",
				Items:       &jsonschema.Schema{Type: "string"},
				Description: "Legacy format: command and arguments as array",
			},
			objectSchema,
		},
	}
}

// UnmarshalYAML implements custom unmarshaling for PostRendererReference.
// Supports both legacy array format and new object format.
func (p *PostRendererReference) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.SequenceNode:
		// Legacy format: post_renderer: ["cmd", "arg1", "arg2"]
		var args []string
		if err := node.Decode(&args); err != nil {
			return fmt.Errorf("failed to decode post_renderer array: %w", err)
		}
		p.Exec = args

		return nil

	case yaml.MappingNode:
		// New object format: post_renderer: { exec: ... } or { kustomize: ... }
		type raw PostRendererReference
		if err := node.Decode((*raw)(p)); err != nil {
			return fmt.Errorf("failed to decode post_renderer object: %w", err)
		}

		// Validate mutual exclusivity
		if len(p.Exec) > 0 && p.Kustomize != nil {
			return fmt.Errorf("post_renderer: exec and kustomize are mutually exclusive")
		}

		return nil

	default:
		return fmt.Errorf("post_renderer must be an array or object, got %v", node.Kind)
	}
}

// UnmarshalYAML implements custom unmarshaling for ExecConfig.
// Supports both string (shlex split) and array formats.
func (e *ExecConfig) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		// String format: exec: "./script.sh --arg1 --arg2"
		var cmd string
		if err := node.Decode(&cmd); err != nil {
			return fmt.Errorf("failed to decode exec string: %w", err)
		}
		args, err := shlex.Split(cmd)
		if err != nil {
			return fmt.Errorf("failed to parse exec command %q: %w", cmd, err)
		}
		*e = args

		return nil

	case yaml.SequenceNode:
		// Array format: exec: ["./script.sh", "--arg1", "--arg2"]
		var args []string
		if err := node.Decode(&args); err != nil {
			return fmt.Errorf("failed to decode exec array: %w", err)
		}
		*e = args

		return nil

	default:
		return fmt.Errorf("exec must be a string or array, got %v", node.Kind)
	}
}

// UnmarshalYAML implements custom unmarshaling for KustomizeConfig.
// Supports both string (directory path) and object formats.
func (k *KustomizeConfig) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		// String format: kustomize: "path/to/dir"
		return node.Decode(&k.Src)

	case yaml.MappingNode:
		// Object format: kustomize: { dir: "path/to/dir", renderer: "sprig" }
		type raw KustomizeConfig

		return node.Decode((*raw)(k))

	default:
		return fmt.Errorf("kustomize must be a string or object, got %v", node.Kind)
	}
}

// IsEmpty returns true if no post-renderer is configured.
func (p *PostRendererReference) IsEmpty() bool {
	return p == nil || (len(p.Exec) == 0 && p.Kustomize == nil)
}

func (p *PostRendererReference) HelmPostRenderer() (postrender.PostRenderer, error) {
	if p.IsEmpty() {
		return nil, nil //nolint:nilnil
	}

	if len(p.Exec) > 0 {
		return postrender.NewExec(p.Exec[0], p.Exec[1:]...) //nolint:wrapcheck
	}

	if p.Kustomize != nil {
		return newKustomizeRenderer(p.Kustomize), nil
	}

	return nil, nil //nolint:nilnil
}

// SetUniq generates unique file path based on provided base directory and release uniqname.
func (k *KustomizeConfig) SetUniq(dir string, name uniqname.UniqName) {
	k.Dst = filepath.Join(dir, "kustomize", name.String())
}

// Build renders the kustomize directory to tmpDir.
// Must be called before PostRenderer() if kustomize mode is used.
func (k *KustomizeConfig) Build(
	ctx context.Context,
	rel Config,
	tmpDir, templater string,
	templateFuncs gotemplate.FuncMap,
) error {
	log.Info("🔨 Building release kustomize...")

	l := rel.Logger().WithField("kustomize src", k.Src)

	// Check source directory exists
	if !helper.IsExists(k.Src) {
		return fmt.Errorf("kustomize directory %q does not exist", k.Src)
	}

	// Generate unique destination path
	k.Dst = filepath.Join(tmpDir, "kustomize", rel.Uniq().String())

	l = l.WithField("kustomize dst", k.Dst)

	// Copy directory to tmpDir
	if err := cp.Copy(k.Src, k.Dst); err != nil {
		return fmt.Errorf("failed to copy kustomize directory %q to %q: %w", k.Src, k.Dst, err)
	}

	// Determine renderer
	renderer := k.Renderer
	if renderer == "" {
		renderer = templater
	}

	// Set default delimiters for gomplate
	delimLeft := k.DelimiterLeft
	delimRight := k.DelimiterRight
	if renderer == template.TemplaterGomplate && delimLeft == "" && delimRight == "" {
		delimLeft = "[["
		delimRight = "]]"
	}

	// Template data
	data := struct {
		Release Config
	}{
		Release: rel,
	}

	opts := []template.TemplaterOptions{
		template.SetDelimiters(delimLeft, delimRight),
	}
	for name, value := range templateFuncs {
		opts = append(opts, template.AddFunc(name, value))
	}

	// Render all files in the directory
	err := filepath.Walk(k.Dst, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		l.WithField("file", path).Trace("Rendering kustomize file")

		return template.Tpl2yml(ctx, path, path, data, renderer, opts...)
	})
	if err != nil {
		return fmt.Errorf("failed to render kustomize directory: %w", err)
	}

	return nil
}

// kustomizeRenderer implements helm's PostRenderer interface using kustomize.
type kustomizeRenderer struct {
	config *KustomizeConfig
}

// newKustomizeRenderer creates a new kustomize post-renderer.
func newKustomizeRenderer(config *KustomizeConfig) *kustomizeRenderer {
	return &kustomizeRenderer{
		config: config,
	}
}

// Run implements postrender.PostRenderer interface.
// It writes helm output to all.yaml, injects it into kustomization and runs kustomize build.
func (r *kustomizeRenderer) Run(renderedManifests *bytes.Buffer) (*bytes.Buffer, error) {
	fSys := filesys.MakeFsOnDisk()

	// Save helm manifests and inject into kustomization
	if err := r.AddHelmManifests(renderedManifests, fSys); err != nil {
		return nil, fmt.Errorf("failed to add helm manifests to kustomization: %w", err)
	}

	// Run kustomize build using krusty API
	kustomizer := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
	resMap, err := kustomizer.Run(fSys, r.config.Dst)
	if err != nil {
		return nil, fmt.Errorf("kustomize build failed in %s: %w", r.config.Dst, err)
	}

	yamlBytes, err := resMap.AsYaml()
	if err != nil {
		return nil, fmt.Errorf("failed to convert kustomize output to YAML: %w", err)
	}

	log.WithField("dir", r.config.Dst).Debug("Kustomize build completed")

	return bytes.NewBuffer(yamlBytes), nil
}

// AddHelmManifests performs a `kustomize edit add resource all.yaml` operation to inject helm manifests into kustomization.
func (r *kustomizeRenderer) AddHelmManifests(renderedManifests *bytes.Buffer, fSys filesys.FileSystem) error {
	kustomizeRoot := filesys.ConfirmedDir(r.config.Dst)
	var kustomizeFile string
	for _, kfilename := range konfig.RecognizedKustomizationFileNames() {
		if fSys.Exists(kustomizeRoot.Join(kfilename)) {
			kustomizeFile = kustomizeRoot.Join(kfilename)

			break
		}
	}

	// Read kustomization file
	data, err := fSys.ReadFile(kustomizeFile)
	if err != nil {
		return err
	}
	var k types.Kustomization
	if err := k.Unmarshal(data); err != nil {
		return err
	}
	k.FixKustomization()

	// Inject helm manifests into kustomization config if not present
	var allYamlFile string
	switch {
	case slices.Index(k.Resources, "all.yaml") != -1:
		allYamlFile = "all.yaml"
	case slices.Index(k.Resources, "all.yml") != -1:
		allYamlFile = "all.yml"
	default:
		allYamlFile = fmt.Sprintf("all%s", filepath.Ext(kustomizeFile))
		k.Resources = append([]string{allYamlFile}, k.Resources...)
		data, err = yaml.Marshal(&k)
		if err != nil {
			return err
		}
		if err := fSys.WriteFile(kustomizeFile, data); err != nil {
			return err
		}
	}

	// Write helm output to all.yaml/all.yml (default to kustomization file extension)
	return os.WriteFile(kustomizeRoot.Join(allYamlFile), renderedManifests.Bytes(), os.FileMode(0o644))
}

func (rel *config) BuildPostRenderer(ctx context.Context, tmpDir, templater string, templateFuncs gotemplate.FuncMap) error {
	if rel.PostRendererF == nil {
		return nil
	}

	if rel.PostRendererF.Kustomize != nil {
		return rel.PostRendererF.Kustomize.Build(ctx, rel, tmpDir, templater, templateFuncs)
	}

	return nil
}
