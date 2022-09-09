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

type config struct {
	cfg                      *action.Configuration
	helm                     *helm.EnvSettings
	log                      *log.Entry
	Store                    map[string]interface{} `json:"store,omitempty" jsonschema:"title=The Store,description=It allows to pass your custom fields from helmwave.yml to values"`
	ChartF                   Chart                  `json:"chart" jsonschema:"oneof_type=string;object"`
	uniqName                 uniqname.UniqName
	NameF                    string            `json:"name" jsonschema:"title=release name"`
	NamespaceF               string            `json:"namespace" jsonschema:"title=kubernetes namespace"`
	DescriptionF             string            `json:"description,omitempty"`
	PendingReleaseStrategy   PendingStrategy   `json:"pending_release_strategy,omitempty" jsonschema:"description=Strategy to handle releases in pending statuses (pending-install, pending-upgrade, pending-rollback)"`
	DependsOnF               []string          `json:"depends_on,omitempty" jsonschema:"title=Needs,description=dependencies"`
	ValuesF                  []ValuesReference `json:"values,omitempty" jsonschema:"title=values of a release"`
	TagsF                    []string          `json:"tags,omitempty" jsonschema:"description=tags allows you choose releases for build"`
	Timeout                  time.Duration     `json:"timeout,omitempty"`
	MaxHistory               int               `json:"max_history,omitempty"`
	AllowFailureF            bool              `json:"allow_failure,omitempty"`
	Atomic                   bool              `json:"atomic,omitempty"`
	CleanupOnFail            bool              `json:"cleanup_on_fail,omitempty"`
	CreateNamespace          bool              `json:"create_namespace,omitempty" jsonschema:"description=will create namespace if it doesnt exits,default=false"`
	Devel                    bool              `json:"devel,omitempty"`
	DisableHooks             bool              `json:"disable_hooks,omitempty"`
	DisableOpenAPIValidation bool              `json:"disable_open_api_validation,omitempty"`
	dryRun                   bool
	Force                    bool `json:"force,omitempty"`
	Recreate                 bool `json:"recreate,omitempty"`
	ResetValues              bool `json:"reset_values,omitempty"`
	ReuseValues              bool `json:"reuse_values,omitempty"`
	SkipCRDs                 bool `json:"skip_crds,omitempty"`
	SubNotes                 bool `json:"sub_notes,omitempty"`
	Wait                     bool `json:"wait,omitempty" jsonschema:"description=prefer use true"`
	WaitForJobs              bool `json:"wait_for_jobs,omitempty" jsonschema:"description=prefer use true"`
}

func (rel *config) DryRun(b bool) {
	rel.dryRun = b
}

// Chart is structure for chart download options.
type Chart struct {
	// action.ChartPathOptions `json:",inline"`
	Version string `json:"version,omitempty" jsonschema:"title=the version,description=The version of a chart,example=0.1.0"`
	Name    string `json:"name" jsonschema:"title=the name,description=The name of a chart,example=bitnami/nginx,example=oci://ghcr.io/helmwave/unit-test-oci"`
}

// UnmarshalYAML flexible Release.
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

	ch := rel.Chart()

	// I hate private field without normal New(...Options)
	// client.ChartPathOptions = ch.ChartPathOptions
	client.ChartPathOptions.Version = ch.Version

	client.DisableHooks = rel.DisableHooks
	client.SkipCRDs = rel.SkipCRDs
	client.Timeout = rel.Timeout
	client.Wait = rel.Wait
	client.WaitForJobs = rel.WaitForJobs
	client.Atomic = rel.Atomic
	client.DisableOpenAPIValidation = rel.DisableOpenAPIValidation
	client.SubNotes = rel.SubNotes
	client.Description = rel.Description()

	if client.DryRun {
		client.Replace = true
		client.ClientOnly = true
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

	ch := rel.Chart()

	// I hate private field without normal New(...Options)
	// client.ChartPathOptions = ch.ChartPathOptions
	client.ChartPathOptions.Version = ch.Version

	client.DisableHooks = rel.DisableHooks
	client.SkipCRDs = rel.SkipCRDs
	client.Timeout = rel.Timeout
	client.Wait = rel.Wait
	client.WaitForJobs = rel.WaitForJobs
	client.Atomic = rel.Atomic
	client.DisableOpenAPIValidation = rel.DisableOpenAPIValidation
	client.SubNotes = rel.SubNotes
	client.Description = rel.Description()

	return client
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
