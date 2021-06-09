package release

import (
	"errors"
	"github.com/helmwave/helmwave/pkg/pubsub"
	"helm.sh/helm/v3/pkg/action"
)

type Config struct {
	Name         string
	Chart        string
	Tags         []string
	Store        map[string]interface{}
	Values       []*ValuesReference
	Options      action.Upgrade
	DependsOn    []string `yaml:"depends_on"`
	dependencies map[string]<-chan pubsub.ReleaseStatus
}

var (
	NotFound = errors.New("release not found")
)
