package release

import (
	"errors"
	"time"

	"github.com/helmwave/helmwave/pkg/pubsub"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/storage/driver"
)

type config struct {
	cfg                      *action.Configuration                             `yaml:"-"`
	dependencies             map[uniqname.UniqName]<-chan pubsub.ReleaseStatus `yaml:"-"`
	helm                     *helm.EnvSettings                                 `yaml:"-"`
	log                      *log.Entry                                        `yaml:"-"`
	Store                    map[string]interface{}                            `yaml:"store,omitempty"`
	ChartF                   Chart                                             `yaml:"chart,omitempty"`
	uniqName                 uniqname.UniqName                                 `yaml:"-"`
	NameF                    string                                            `yaml:"name,omitempty"`
	NamespaceF               string                                            `yaml:"namespace,omitempty"`
	DescriptionF             string                                            `yaml:"description,omitempty"`
	DependsOnF               []string                                          `yaml:"depends_on,omitempty"`
	ValuesF                  []ValuesReference                                 `yaml:"values,omitempty"`
	TagsF                    []string                                          `yaml:"tags,omitempty"`
	Timeout                  time.Duration                                     `yaml:"timeout,omitempty"`
	MaxHistory               int                                               `yaml:"maxhistory,omitempty"`
	AllowFailure             bool                                              `yaml:"allow_failure,omitempty"`
	Atomic                   bool                                              `yaml:"atomic,omitempty"`
	CleanupOnFail            bool                                              `yaml:"cleanuponfail,omitempty"`
	CreateNamespace          bool                                              `yaml:"createnamespace,omitempty"`
	Devel                    bool                                              `yaml:"devel,omitempty"`
	DisableHooks             bool                                              `yaml:"disablehooks,omitempty"`
	DisableOpenAPIValidation bool                                              `yaml:"disableopenapivalidation,omitempty"`
	dryRun                   bool                                              `yaml:"dryrun,omitempty"`
	Force                    bool                                              `yaml:"force,omitempty"`
	Recreate                 bool                                              `yaml:"recreate,omitempty"`
	ResetValues              bool                                              `yaml:"resetvalues,omitempty"`
	ReuseValues              bool                                              `yaml:"reusevalues,omitempty"`
	SkipCRDs                 bool                                              `yaml:"skipcrds,omitempty"`
	SubNotes                 bool                                              `yaml:"subnotes,omitempty"`
	Wait                     bool                                              `yaml:"wait,omitempty"`
	WaitForJobs              bool                                              `yaml:"waitforjobs,omitempty"`
}

func (rel *config) DryRun(b bool) {
	rel.dryRun = b
}

// Chart is structure for chart download options.
type Chart struct {
	action.ChartPathOptions `yaml:",inline"`
	Name                    string
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
	client.ChartPathOptions = rel.Chart().ChartPathOptions
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
	client.ChartPathOptions = rel.Chart().ChartPathOptions
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

// In check that 'x' found in 'array'.
func (rel *config) In(a []Config) bool {
	for _, r := range a {
		if rel.Uniq() == r.Uniq() {
			return true
		}
	}

	return false
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

func (rel *config) DependsOn() []string {
	return rel.DependsOnF
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
