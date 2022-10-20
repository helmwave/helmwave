package release

import (
	"errors"
	"fmt"
	"time"

	"github.com/invopop/jsonschema"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/postrender"
	"helm.sh/helm/v3/pkg/storage/driver"
)

// Configs type of array Config.
type Configs []Config

// UnmarshalYAML parse Config.
func (r *Configs) UnmarshalYAML(node *yaml.Node) error {
	var err error
	*r, err = UnmarshalYAML(node)

	return err
}

func (Configs) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{DoNotReference: true}
	var l []*config

	return r.Reflect(&l)
}

//nolint:lll
type config struct {
	helm                     *helm.EnvSettings      `yaml:"-"`
	log                      *log.Entry             `yaml:"-"`
	Store                    map[string]interface{} `yaml:"store,omitempty" json:"store,omitempty" jsonschema:"title=The Store,description=It allows to pass your custom fields from helmwave.yml to values"`
	ChartF                   Chart                  `yaml:"chart,omitempty" json:"chart,omitempty" jsonschema:"oneof_type=string;object"`
	uniqName                 uniqname.UniqName      `yaml:"-"`
	NameF                    string                 `yaml:"name,omitempty" json:"name,omitempty" jsonschema:"title=release name"`
	NamespaceF               string                 `yaml:"namespace,omitempty" json:"namespace,omitempty" jsonschema:"title=kubernetes namespace"`
	DescriptionF             string                 `yaml:"description,omitempty" json:"description,omitempty"`
	PendingReleaseStrategy   PendingStrategy        `yaml:"pending_release_strategy,omitempty" json:"pending_release_strategy,omitempty" jsonschema:"description=Strategy to handle releases in pending statuses (pending-install, pending-upgrade, pending-rollback)"`
	DependsOnF               []string               `yaml:"depends_on,omitempty" json:"depends_on,omitempty" jsonschema:"title=Needs,description=dependencies"`
	ValuesF                  []ValuesReference      `yaml:"values,omitempty" json:"values,omitempty" jsonschema:"title=values of a release"`
	TagsF                    []string               `yaml:"tags,omitempty" json:"tags,omitempty" jsonschema:"description=tags allows you choose releases for build"`
	PostRendererF            []string               `yaml:"post_renderer,omitempty" json:"post_renderer,omitempty"`
	Timeout                  time.Duration          `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	MaxHistory               int                    `yaml:"max_history,omitempty" json:"max_history,omitempty"`
	AllowFailureF            bool                   `yaml:"allow_failure,omitempty" json:"allow_failure,omitempty"`
	Atomic                   bool                   `yaml:"atomic,omitempty" json:"atomic,omitempty"`
	CleanupOnFail            bool                   `yaml:"cleanup_on_fail,omitempty" json:"cleanup_on_fail,omitempty"`
	CreateNamespace          bool                   `yaml:"create_namespace,omitempty" json:"create_namespace,omitempty" jsonschema:"description=will create namespace if it doesnt exits,default=false"`
	Devel                    bool                   `yaml:"devel,omitempty" json:"devel,omitempty"`
	DisableHooks             bool                   `yaml:"disable_hooks,omitempty" json:"disable_hooks,omitempty"`
	DisableOpenAPIValidation bool                   `yaml:"disable_open_api_validation,omitempty" json:"disable_open_api_validation,omitempty"`
	dryRun                   bool                   `yaml:"dry_run,omitempty" json:"dry_run,omitempty"` //nolint:govet
	Force                    bool                   `yaml:"force,omitempty" json:"force,omitempty"`
	Recreate                 bool                   `yaml:"recreate,omitempty" json:"recreate,omitempty"`
	ResetValues              bool                   `yaml:"reset_values,omitempty" json:"reset_values,omitempty"`
	ReuseValues              bool                   `yaml:"reuse_values,omitempty" json:"reuse_values,omitempty"`
	SkipCRDs                 bool                   `yaml:"skip_crds,omitempty" json:"skip_crds,omitempty"`
	SubNotes                 bool                   `yaml:"sub_notes,omitempty" json:"sub_notes,omitempty"`
	Wait                     bool                   `yaml:"wait,omitempty" json:"wait,omitempty" jsonschema:"description=prefer use true"`
	WaitForJobs              bool                   `yaml:"wait_for_jobs,omitempty" json:"wait_for_jobs,omitempty" jsonschema:"description=prefer use true"`
}

func (rel *config) DryRun(b bool) {
	rel.dryRun = b
}

// Chart is structure for chart download options.
//
//nolint:lll
type Chart struct {
	action.ChartPathOptions `yaml:",inline"`
	Name                    string `yaml:"name" json:"name" jsonschema:"title=the name,description=The name of a chart,example=bitnami/nginx,example=oci://ghcr.io/helmwave/unit-test-oci"`
}

// UnmarshalYAML flexible config.
func (u *Chart) UnmarshalYAML(node *yaml.Node) error {
	type raw Chart
	var err error

	switch node.Kind {
	case yaml.ScalarNode, yaml.AliasNode:
		err = node.Decode(&(u.Name))
	case yaml.MappingNode:
		err = node.Decode((*raw)(u))
	default:
		err = fmt.Errorf("unknown format")
	}

	if err != nil {
		return fmt.Errorf("failed to decode chart %q from YAML at %d line: %w", node.Value, node.Line, err)
	}

	return nil
}

func (rel *config) newInstall() *action.Install {
	client := action.NewInstall(rel.Cfg())

	// Only Up
	client.CreateNamespace = rel.CreateNamespace
	client.ReleaseName = rel.Name()

	// Common Part
	client.DryRun = rel.dryRun
	client.Devel = rel.Devel
	client.Namespace = rel.Namespace()

	rel.copyChartPathOptions(&client.ChartPathOptions)

	client.DisableHooks = rel.DisableHooks
	client.SkipCRDs = rel.SkipCRDs
	client.Timeout = rel.Timeout
	client.Wait = rel.Wait
	client.WaitForJobs = rel.WaitForJobs
	client.Atomic = rel.Atomic
	client.DisableOpenAPIValidation = rel.DisableOpenAPIValidation
	client.SubNotes = rel.SubNotes
	client.Description = rel.Description()

	pr, err := rel.PostRenderer()
	if err != nil {
		rel.Logger().WithError(err).Warn("failed to create postrenderer")
	} else {
		client.PostRenderer = pr
	}

	if client.DryRun {
		client.Replace = true
	}

	return client
}

func (rel *config) newUpgrade() *action.Upgrade {
	client := action.NewUpgrade(rel.Cfg())

	// Only Upgrade
	client.CleanupOnFail = rel.CleanupOnFail
	client.MaxHistory = rel.MaxHistory
	client.Recreate = rel.Recreate
	client.ReuseValues = rel.ReuseValues
	client.ResetValues = rel.ResetValues

	// Common Part
	client.DryRun = rel.dryRun
	client.Devel = rel.Devel
	client.Namespace = rel.Namespace()

	rel.copyChartPathOptions(&client.ChartPathOptions)

	client.Force = rel.Force
	client.DisableHooks = rel.DisableHooks
	client.SkipCRDs = rel.SkipCRDs
	client.Timeout = rel.Timeout
	client.Wait = rel.Wait
	client.WaitForJobs = rel.WaitForJobs
	client.Atomic = rel.Atomic
	client.DisableOpenAPIValidation = rel.DisableOpenAPIValidation
	client.SubNotes = rel.SubNotes
	client.Description = rel.Description()

	pr, err := rel.PostRenderer()
	if err != nil {
		rel.Logger().WithError(err).Warn("failed to create postrenderer")
	} else {
		client.PostRenderer = pr
	}

	return client
}

func (rel *config) copyChartPathOptions(cpo *action.ChartPathOptions) {
	ch := rel.Chart()

	// I hate private field without normal New(...Options)
	cpo.CaFile = ch.ChartPathOptions.CaFile
	cpo.CertFile = ch.ChartPathOptions.CertFile
	cpo.KeyFile = ch.ChartPathOptions.KeyFile
	cpo.InsecureSkipTLSverify = ch.ChartPathOptions.InsecureSkipTLSverify
	cpo.Keyring = ch.ChartPathOptions.Keyring
	cpo.Password = ch.ChartPathOptions.Password
	cpo.PassCredentialsAll = ch.ChartPathOptions.PassCredentialsAll
	cpo.RepoURL = ch.ChartPathOptions.RepoURL
	cpo.Username = ch.ChartPathOptions.Username
	cpo.Verify = ch.ChartPathOptions.Verify
	cpo.Version = ch.ChartPathOptions.Version
}

var (
	// ErrNotFound is an error for not found release.
	ErrNotFound = driver.ErrReleaseNotFound

	// ErrFoundMultiple is an error for multiple releases found by name.
	ErrFoundMultiple = errors.New("found multiple releases o_0")

	// ErrDepFailed is an error thrown when dependency release fails.
	ErrDepFailed = errors.New("dependency failed")
)

// Uniq redis@my-namespace.
func (rel *config) Uniq() uniqname.UniqName {
	if rel.uniqName == "" {
		var err error
		rel.uniqName, err = uniqname.Generate(rel.Name(), rel.Namespace())
		if err != nil {
			rel.Logger().WithFields(log.Fields{
				"name":       rel.Name(),
				"namespace":  rel.Namespace(),
				log.ErrorKey: err,
			}).Error("failed to generate valid uniqname")
		}
	}

	return rel.uniqName
}

func (rel *config) Equal(a Config) bool {
	return rel.Uniq().Equal(a.Uniq())
}

func (rel *config) Name() string {
	return rel.NameF
}

func (rel *config) Namespace() string {
	return rel.NamespaceF
}

func (rel *config) Description() string {
	return rel.DescriptionF
}

func (rel *config) Chart() Chart {
	return rel.ChartF
}

func (rel *config) DependsOn() []uniqname.UniqName {
	result := make([]uniqname.UniqName, len(rel.DependsOnF))

	for i, dep := range rel.DependsOnF {
		result[i] = uniqname.UniqName(dep)
	}

	return result
}

func (rel *config) Tags() []string {
	return rel.TagsF
}

func (rel *config) Values() []ValuesReference {
	return rel.ValuesF
}

func (rel *config) Logger() *log.Entry {
	if rel.log == nil {
		rel.log = log.WithField("release", rel.Uniq())
	}

	return rel.log
}

func (rel *config) AllowFailure() bool {
	return rel.AllowFailureF
}

func (rel *config) HelmWait() bool {
	return rel.Wait
}

func (rel *config) buildAfterUnmarshal() {
	rel.buildAfterUnmarshalDependsOn()

	// set default timeout
	if rel.Timeout <= 0 {
		rel.Timeout = 5 * time.Minute
	}
}

func (rel *config) buildAfterUnmarshalDependsOn() {
	res := make([]string, 0, len(rel.DependsOnF))

	for _, dep := range rel.DependsOnF {
		u, err := uniqname.GenerateWithDefaultNamespace(dep, rel.Namespace())
		if err != nil {
			rel.Logger().WithError(err).WithField("dependency", dep).Error("Cannot parse dependency")

			continue
		}

		// generate full uniqname string if it was short
		res = append(res, string(u))
	}

	rel.DependsOnF = res
}

func (rel *config) PostRenderer() (postrender.PostRenderer, error) {
	if len(rel.PostRendererF) < 1 {
		return nil, nil
	}

	return postrender.NewExec(rel.PostRendererF[0], rel.PostRendererF[1:]...) //nolint:wrapcheck
}
