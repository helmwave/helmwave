package release

import (
	"github.com/helmwave/helmwave/pkg/pubsub"
	"helm.sh/helm/v3/pkg/action"
)

type Config struct {
	dependencies map[string]<-chan pubsub.ReleaseStatus
	Store        map[string]interface{}
	Options      action.Upgrade
	Name         string
	Chart        string
	Tags         []string
	Values       []string
	DependsOn    []string `yaml:"depends_on"`
}
