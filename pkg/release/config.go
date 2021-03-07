package release

import (
	"github.com/zhilyaev/helmwave/pkg/pubsub"
	"helm.sh/helm/v3/pkg/action"
)

type Config struct {
	Name         string
	Chart        string
	Tags         []string
	Store        map[string]interface{}
	Values       []string
	Options      action.Upgrade
	DependsOn    []string `json:"depends_on"`
	dependencies map[string]<-chan pubsub.ReleaseStatus
}
