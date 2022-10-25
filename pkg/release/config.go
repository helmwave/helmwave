package release

import (
	"errors"
	"sync"
	"time"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/postrender"
	"helm.sh/helm/v3/pkg/storage/driver"
)

//nolint:lll
type config struct {
	helm                     *helm.EnvSettings      `json:"-"`
	log                      *log.Entry             `json:"-"`
	Store                    map[string]interface{} `json:"store,omitempty" jsonschema:"title=The Store,description=It allows to pass your custom fields from helmwave.yml to values"`
	ChartF                   Chart                  `json:"chart,omitempty" jsonschema:"title=Chart reference,description=Describes chart that release uses,oneof_type=string;object"`
	PendingReleaseStrategy   PendingStrategy        `json:"pending_release_strategy,omitempty" jsonschema:"description=Strategy to handle releases in pending statuses (pending-install/pending-upgrade/pending-rollback),default="`
	uniqName                 uniqname.UniqName      `json:"-"`
	NameF                    string                 `json:"name,omitempty" jsonschema:"required,title=Release name"`
	NamespaceF               string                 `json:"namespace,omitempty" jsonschema:"required,title=Kubernetes namespace"`
	DescriptionF             string                 `json:"description,omitempty" jsonschema:"default="`
	KubeContextF             string                 `json:"context,omitempty"`
	DependsOnF               []*DependsOnReference  `json:"depends_on,omitempty" jsonschema:"title=Needs,description=List of dependencies that are required to succeed before this release"`
	ValuesF                  []ValuesReference      `json:"values,omitempty" jsonschema:"title=Values of the release"`
	TagsF                    []string               `json:"tags,omitempty" jsonschema:"description=Tags allows you choose releases for build"`
	PostRendererF            []string               `json:"post_renderer,omitempty" jsonschema:"description=List of postrenders to manipulate with manifests"`
	MaxHistory               int                    `json:"max_history,omitempty" jsonschema:"default=0"`
	Timeout                  time.Duration          `json:"timeout,omitempty" jsonschema:"default=5m"`
	lock                     sync.RWMutex           `json:"-"`
	AllowFailureF            bool                   `json:"allow_failure,omitempty" jsonschema:"description=Whether to ignore errors and proceed with dependant releases,default=false"`
	Atomic                   bool                   `json:"atomic,omitempty" jsonschema:"default=false"`
	CleanupOnFail            bool                   `json:"cleanup_on_fail,omitempty" jsonschema:"default=false"`
	CreateNamespace          bool                   `json:"create_namespace,omitempty" jsonschema:"description=Whether to create namespace if it doesnt exits,default=false"`
	Devel                    bool                   `json:"devel,omitempty" jsonschema:"default=false"`
	DisableHooks             bool                   `json:"disable_hooks,omitempty" jsonschema:"default=false"`
	DisableOpenAPIValidation bool                   `json:"disable_open_api_validation,omitempty" jsonschema:"default=false"`
	dryRun                   bool                   `json:"dry_run,omitempty" jsonschema:"default=false"` //nolint:govet
	Force                    bool                   `json:"force,omitempty" jsonschema:"default=false"`
	Recreate                 bool                   `json:"recreate,omitempty" jsonschema:"default=false"`
	ResetValues              bool                   `json:"reset_values,omitempty" jsonschema:"default=false"`
	ReuseValues              bool                   `json:"reuse_values,omitempty" jsonschema:"default=false"`
	SkipCRDs                 bool                   `json:"skip_crds,omitempty" jsonschema:"default=false"`
	SubNotes                 bool                   `json:"sub_notes,omitempty" jsonschema:"default=false"`
	Wait                     bool                   `json:"wait,omitempty" jsonschema:"description=Whether to wait for all resource to become ready,default=false"`
	WaitForJobs              bool                   `json:"wait_for_jobs,omitempty" jsonschema:"description=Whether to wait for all jobs to become ready,default=false"`
}

func (rel *config) DryRun(b bool) {
	rel.dryRun = b
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
	rel.lock.RLock()
	defer rel.lock.RUnlock()

	return rel.ChartF
}

func (rel *config) DependsOn() []*DependsOnReference {
	rel.lock.RLock()
	defer rel.lock.RUnlock()

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
		rel.Logger().WithField("dependency", dep).WithError(err).Error("Cannot parse dependency")

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

// MarshalYAML is a marshaller for github.com/goccy/go-yaml.
// It is required to avoid data race with getting read lock.
//
//nolint:govet
func (rel *config) MarshalYAML() (interface{}, error) {
	rel.lock.RLock()
	defer rel.lock.RUnlock()

	type raw config
	r := raw(*rel)

	return r, nil
}
