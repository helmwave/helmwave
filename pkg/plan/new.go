package plan

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/registry"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/pkg/version"
	"github.com/invopop/jsonschema"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	// Dir is default directory for generated files.
	Dir = ".helmwave/"

	// File is default file name for planfile.
	File = "planfile"

	// Body is default file name for main config.
	Body = "helmwave.yml"

	// Manifest is default directory under Dir for manifests.
	Manifest = "manifest/"

	// Values is default directory for values.
	Values = "values/"
)

var (
	// ErrManifestDirNotFound is an error for nonexistent manifest dir.
	ErrManifestDirNotFound = errors.New(Manifest + " dir not found")

	// ErrManifestDirEmpty is an error for empty manifest dir.
	ErrManifestDirEmpty = errors.New(Manifest + " is empty")
)

// Plan contains full helmwave state.
type Plan struct {
	body      *planBody
	dir       string
	fullPath  string
	tmpDir    string
	graphMD   string
	templater string

	manifests map[uniqname.UniqName]string
	unchanged release.Configs
}

// NewAndImport wrapper for New and Import in one.
func NewAndImport(ctx context.Context, src string) (p *Plan, err error) {
	p = New(src)

	err = p.Import(ctx)
	if err != nil {
		return p, err
	}

	return p, nil
}

// Logger will pretty build log.Entry.
func (p *Plan) Logger() *log.Entry {
	a := make([]string, 0, len(p.body.Releases))
	for _, r := range p.body.Releases {
		a = append(a, r.Uniq().String())
	}

	b := make([]string, 0, len(p.body.Repositories))
	for _, r := range p.body.Repositories {
		b = append(b, r.Name())
	}

	c := make([]string, 0, len(p.body.Registries))
	for _, r := range p.body.Registries {
		c = append(c, r.Host())
	}

	return log.WithFields(log.Fields{
		"releases":     a,
		"repositories": b,
		"registries":   c,
	})
}

// nolintlint:lll
type planBody struct {
	Project      string           `yaml:"project" json:"project" jsonschema:"title=project name,description=reserved for future,example=my-awesome-project"`                                                                         //nolint:lll
	Version      string           `yaml:"version" json:"version" jsonschema:"title=version of helmwave,description=will check current version and project version,pattern=^[0-9]+\\.[0-9]+\\.[0-9]+$,example=0.23.0,example=0.22.1"` //nolint:lll
	Repositories repo.Configs     `yaml:"repositories" json:"repositories" jsonschema:"title=repositories list,description=helm repositories"`                                                                                       //nolint:lll
	Registries   registry.Configs `yaml:"registries" json:"registries" jsonschema:"title=registries list,description=helm OCI registries"`                                                                                           //nolint:lll
	Releases     release.Configs  `yaml:"releases" json:"releases" jsonschema:"title=helm releases,description=what you wanna deploy"`                                                                                               //nolint:lll
}

func GenSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
	}

	schema := r.Reflect(&planBody{})
	schema.AdditionalProperties = jsonschema.TrueSchema // to allow anchors at the top level

	return schema
}

// NewBody parses plan from file.
func NewBody(ctx context.Context, file string) (*planBody, error) {
	b := &planBody{
		Version: version.Version,
	}

	src, err := os.ReadFile(file)
	if err != nil {
		return b, fmt.Errorf("failed to read plan file %s: %w", file, err)
	}

	decoder := yaml.NewDecoder(bytes.NewBuffer(src))
	err = decoder.Decode(b)
	if err != nil {
		return b, fmt.Errorf("failed to unmarshal YAML plan %s: %w", file, err)
	}

	// Setup dev version
	// if b.Version == "" {
	// 	 b.Version = version.Version
	// }

	err = b.Validate()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// New returns empty *Plan for provided directory.
func New(dir string) *Plan {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		log.WithError(err).Warn("failed to create temporary directory")
		tmpDir = os.TempDir()
	}

	plan := &Plan{
		tmpDir:    tmpDir,
		dir:       dir,
		fullPath:  filepath.Join(dir, File),
		manifests: make(map[uniqname.UniqName]string),
	}

	return plan
}
