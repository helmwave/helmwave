package release

import (
	"errors"
	"github.com/helmwave/helmwave/pkg/pubsub"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/storage/driver"
)

type Config struct {
	action.Install `yaml:",inline"`

	//action.Upgrade
	Force, ResetValues, ReuseValues, Recreate, CleanupOnFail bool
	MaxHistory int


	// Helmwave
	Chart        string
	Tags         []string
	Values       []string
	Store        map[string]interface{}
	DependsOn    []string `yaml:"depends_on"`
	dependencies map[string]<-chan pubsub.ReleaseStatus
}

var (
	ErrNotFound      = driver.ErrReleaseNotFound
	ErrFoundMultiple = errors.New("found multiple releases o_0")
	ErrEmpty         = errors.New("releases are empty")
	ErrDepFailed 	 = errors.New("dependency failed")
)

// UniqName redis@my-namespace
func (rel *Config) UniqName() string {
	return rel.ReleaseName + "@" + rel.Namespace
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

