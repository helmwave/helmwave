package release

import (
	"errors"
	"time"

	"github.com/helmwave/helmwave/pkg/pubsub"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/storage/driver"
)

type Config struct {
	cfg                      *action.Configuration
	dependencies             map[uniqname.UniqName]<-chan pubsub.ReleaseStatus
	helm                     *helm.EnvSettings
	Store                    map[string]interface{}
	Chart                    Chart
	uniqName                 uniqname.UniqName
	Name                     string            `yaml:"name"`
	Namespace                string            `yaml:"namespace"`
	Description              string            `yaml:"description"`
	DependsOn                []string          `yaml:"depends_on"`
	Values                   []ValuesReference `yaml:"values"`
	Tags                     []string          `yaml:"tags"`
	Timeout                  time.Duration     `yaml:"timeout"`
	MaxHistory               int
	AllowFailure             bool `yaml:"allow_failure"`
	CreateNamespace          bool
	ResetValues              bool
	Recreate                 bool
	Force                    bool
	Atomic                   bool
	CleanupOnFail            bool
	SubNotes                 bool
	DisableHooks             bool
	DisableOpenAPIValidation bool
	WaitForJobs              bool
	Wait                     bool
	SkipCRDs                 bool
	dryRun                   bool
	Devel                    bool
	ReuseValues              bool
}

func (rel *Config) DryRun(b bool) *Config {
	rel.dryRun = b
	return rel
}

type Chart struct {
	Name                    string
	action.ChartPathOptions `yaml:",inline"`
}

func (rel *Config) newInstall() *action.Install {
	client := action.NewInstall(rel.Cfg())

	// Only Up
	client.CreateNamespace = rel.CreateNamespace
	client.ReleaseName = rel.Name

	// Common Part
	client.DryRun = rel.dryRun
	client.Devel = rel.Devel
	client.Namespace = rel.Namespace
	client.ChartPathOptions = rel.Chart.ChartPathOptions
	client.DisableHooks = rel.DisableHooks
	client.SkipCRDs = rel.SkipCRDs
	client.Timeout = rel.Timeout
	client.Wait = rel.Wait
	client.Atomic = rel.Atomic
	client.DisableOpenAPIValidation = rel.DisableOpenAPIValidation
	client.SubNotes = rel.SubNotes
	client.Description = rel.Description

	if client.DryRun {
		client.Replace = true
		client.ClientOnly = true
	}

	return client
}

func (rel *Config) newUpgrade() *action.Upgrade {
	client := action.NewUpgrade(rel.Cfg())
	// Only Upgrade
	client.CleanupOnFail = rel.CleanupOnFail
	client.MaxHistory = rel.MaxHistory
	client.Recreate = rel.Recreate
	client.ReuseValues = rel.ReuseValues
	client.ResetValues = rel.ReuseValues

	// Common Part
	client.DryRun = rel.dryRun
	client.Devel = rel.Devel
	client.Namespace = rel.Namespace
	client.ChartPathOptions = rel.Chart.ChartPathOptions
	client.DisableHooks = rel.DisableHooks
	client.SkipCRDs = rel.SkipCRDs
	client.Timeout = rel.Timeout
	client.Wait = rel.Wait
	client.Atomic = rel.Atomic
	client.DisableOpenAPIValidation = rel.DisableOpenAPIValidation
	client.SubNotes = rel.SubNotes
	client.Description = rel.Description

	return client
}

var (
	ErrNotFound      = driver.ErrReleaseNotFound
	ErrFoundMultiple = errors.New("found multiple releases o_0")
	ErrDepFailed     = errors.New("dependency failed")
)

// Uniq redis@my-namespace
func (rel *Config) Uniq() uniqname.UniqName {
	if rel.uniqName == "" {
		rel.uniqName = uniqname.UniqName(rel.Name + uniqname.Separator + rel.Namespace)
	}

	return rel.uniqName
}

// In check that 'x' found in 'array'
func (rel *Config) In(a []*Config) bool {
	for _, r := range a {
		if rel == r {
			return true
		}
	}
	return false
}
