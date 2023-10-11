package release

import (
	"sync"
	"time"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/hooks"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/chartutil"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/postrender"
)

type config struct {
	helm *helm.EnvSettings
	log  *log.Entry

	Lifecycle              hooks.Lifecycle `yaml:"lifecycle,omitempty" json:"lifecycle,omitempty" jsonschema:"description=Lifecycle hooks"`
	Store                  map[string]any  `yaml:"store,omitempty" json:"store,omitempty" jsonschema:"title=The Store,description=It allows to pass your custom fields from helmwave.yml to values"`
	ChartF                 Chart           `yaml:"chart,omitempty" json:"chart,omitempty" jsonschema:"title=Chart reference,description=Describes chart that release uses,oneof_type=string;object"`
	PendingReleaseStrategy PendingStrategy `yaml:"pending_release_strategy,omitempty" json:"pending_release_strategy,omitempty" jsonschema:"description=Strategy to handle releases in pending statuses (pending-install/pending-upgrade/pending-rollback)"`
	uniqName               uniqname.UniqName

	NameF               string `yaml:"name,omitempty" json:"name,omitempty" jsonschema:"required,title=Release name"`
	NamespaceF          string `yaml:"namespace,omitempty" json:"namespace,omitempty" jsonschema:"required,title=Kubernetes namespace"`
	DescriptionF        string `yaml:"description,omitempty" json:"description,omitempty" jsonschema:"default="`
	KubeContextF        string `yaml:"context,omitempty" json:"context,omitempty"`
	OfflineKubeVersionF string `yaml:"offline_kube_version,omitempty" json:"offline_kube_version,omitempty" jsonschema:"description=Kubernetes version for offline mode"`

	DependsOnF    []*DependsOnReference `yaml:"depends_on,omitempty" json:"depends_on,omitempty" jsonschema:"title=Needs,description=List of dependencies that are required to succeed before this release"`
	MonitorsF     []MonitorReference    `yaml:"monitors,omitempty" json:"monitors,omitempty" jsonschema:"title=Monitors to execute after upgrade"`
	ValuesF       []ValuesReference     `yaml:"values,omitempty" json:"values,omitempty" jsonschema:"title=Values of the release,oneof_type=string;object"`
	TagsF         []string              `yaml:"tags,omitempty" json:"tags,omitempty" jsonschema:"description=Tags allows you choose releases for build"`
	PostRendererF []string              `yaml:"post_renderer,omitempty" json:"post_renderer,omitempty" jsonschema:"description=List of post_renders to manipulate with manifests"`
	ShowNotes     bool                  `yaml:"show_notes,omitempty" json:"show_notes,omitempty" jsonschema:"description=Output rendered chart notes after upgrade/install"`

	Timeout time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty" jsonschema:"oneof_type=string;int,default=5m"`

	// Lock for parallel testing
	lock sync.RWMutex

	MaxHistory               int  `yaml:"max_history,omitempty" json:"max_history,omitempty" jsonschema:"default=0"`
	AllowFailureF            bool `yaml:"allow_failure,omitempty" json:"allow_failure,omitempty" jsonschema:"description=Whether to ignore errors and proceed with dependant releases,default=false"`
	Atomic                   bool `yaml:"atomic,omitempty" json:"atomic,omitempty" jsonschema:"default=false"`
	CleanupOnFail            bool `yaml:"cleanup_on_fail,omitempty" json:"cleanup_on_fail,omitempty" jsonschema:"default=false"`
	CreateNamespace          bool `yaml:"create_namespace,omitempty" json:"create_namespace,omitempty" jsonschema:"description=Whether to create namespace if it doesnt exits,default=false"`
	DisableHooks             bool `yaml:"disable_hooks,omitempty" json:"disable_hooks,omitempty" jsonschema:"default=false"`
	DisableOpenAPIValidation bool `yaml:"disable_open_api_validation,omitempty" json:"disable_open_api_validation,omitempty" jsonschema:"default=false"`
	EnableDNS                bool `yaml:"enable_dns,omitempty" json:"enable_dns,omitempty" jsonschema:"default=false"`
	Force                    bool `yaml:"force,omitempty" json:"force,omitempty" jsonschema:"default=false"`
	Recreate                 bool `yaml:"recreate,omitempty" json:"recreate,omitempty" jsonschema:"default=false"`
	ResetValues              bool `yaml:"reset_values,omitempty" json:"reset_values,omitempty" jsonschema:"default=false"`
	ReuseValues              bool `yaml:"reuse_values,omitempty" json:"reuse_values,omitempty" jsonschema:"default=false"`
	SkipCRDs                 bool `yaml:"skip_crds,omitempty" json:"skip_crds,omitempty" jsonschema:"default=false"`
	SubNotes                 bool `yaml:"sub_notes,omitempty" json:"sub_notes,omitempty" jsonschema:"default=false"`
	Wait                     bool `yaml:"wait,omitempty" json:"wait,omitempty" jsonschema:"description=Whether to wait for all resource to become ready,default=false"`
	WaitForJobs              bool `yaml:"wait_for_jobs,omitempty" json:"wait_for_jobs,omitempty" jsonschema:"description=Whether to wait for all jobs to become ready,default=false"`

	Labels map[string]string `yaml:"labels,omitempty" json:"labels,omitempty" jsonschema:"Labels that would be added to release metadata on sync"`

	// special field for templating and building
	dryRun bool `jsonschema:"default=false,-"`
}

func (rel *config) DryRun(b bool) {
	rel.dryRun = b
}

// Uniq like redis@my-namespace.
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

func (rel *config) Chart() *Chart {
	rel.lock.RLock()
	defer rel.lock.RUnlock()

	return &rel.ChartF
}

func (rel *config) DependsOn() []*DependsOnReference {
	rel.lock.RLock()
	defer rel.lock.RUnlock()

	return rel.DependsOnF
}

func (rel *config) SetDependsOn(deps []*DependsOnReference) {
	rel.lock.Lock()
	defer rel.lock.Unlock()

	rel.DependsOnF = deps
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

func (rel *config) buildAfterUnmarshal(allReleases []*config) {
	rel.buildAfterUnmarshalDependsOn(allReleases)

	// set default timeout
	if rel.Timeout <= 0 {
		rel.Logger().Debug("timeout is not set, defaulting to 5m")
		rel.Timeout = 5 * time.Minute
	}
}

func (rel *config) buildAfterUnmarshalDependsOn(allReleases []*config) {
	newDeps := make([]*DependsOnReference, 0)

	for _, dep := range rel.DependsOn() {
		l := rel.Logger().WithField("dependency", dep)
		switch dep.Type() {
		case DependencyRelease:
			err := rel.buildAfterUnmarshalDependency(dep)
			if err == nil {
				newDeps = append(newDeps, dep)
			}
		case DependencyTag:
			for _, r := range allReleases {
				if helper.Contains(dep.Tag, r.Tags()) {
					newDep := &DependsOnReference{
						Name:     r.Uniq().String(),
						Optional: dep.Optional,
					}
					newDeps = append(newDeps, newDep)
				}
			}
		case DependencyInvalid:
			l.Warn("invalid dependency, skipping")
		}
	}

	rel.lock.Lock()
	rel.DependsOnF = newDeps
	rel.lock.Unlock()
}

func (rel *config) buildAfterUnmarshalDependency(dep *DependsOnReference) error {
	u, err := uniqname.GenerateWithDefaultNamespace(dep.Name, rel.Namespace())
	if err != nil {
		rel.Logger().WithField("dependency", dep).WithError(err).Error("can't parse dependency")

		return err
	}

	// generate full uniqname string if it was short
	dep.Name = u.String()

	return nil
}

func (rel *config) PostRenderer() (postrender.PostRenderer, error) {
	if len(rel.PostRendererF) < 1 {
		return nil, nil
	}

	return postrender.NewExec(rel.PostRendererF[0], rel.PostRendererF[1:]...) //nolint:wrapcheck
}

func (rel *config) KubeContext() string {
	return rel.KubeContextF
}

// MarshalYAML is a marshaller for gopkg.in/yaml.v3.
// It is required to avoid data race with getting read lock.
func (rel *config) MarshalYAML() (any, error) {
	rel.lock.RLock()
	defer rel.lock.RUnlock()

	type raw config
	r := raw(*rel) //nolint:govet

	return r, nil //nolint:govet
}

func (rel *config) HooksDisabled() bool {
	return rel.DisableHooks
}

func (rel *config) OfflineKubeVersion() *chartutil.KubeVersion {
	if rel.OfflineKubeVersionF != "" {
		v, err := chartutil.ParseKubeVersion(rel.OfflineKubeVersionF)
		if err != nil {
			log.Fatalf("invalid kube version %q: %s", rel.OfflineKubeVersionF, err)

			return nil
		}

		return v
	}

	return nil
}

func (rel *config) Monitors() []MonitorReference {
	rel.lock.RLock()
	defer rel.lock.RUnlock()

	return rel.MonitorsF
}
