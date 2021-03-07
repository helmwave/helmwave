package release

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/zhilyaev/helmwave/pkg/pubsub"
)

var releasePubSub = pubsub.NewReleasePubSub()
var DependencyFailedError = errors.New("dependency failed")

func (rel *Config) NotifySuccess() {
	releasePubSub.PublishSuccess(rel.Name)
}

func (rel *Config) NotifyFailed() {
	releasePubSub.PublishFailed(rel.Name)
}

func (rel *Config) addDependency(name string) {
	ch := releasePubSub.Subscribe(name)

	if rel.dependencies == nil {
		rel.dependencies = make(map[string]<-chan pubsub.ReleaseStatus)
	}

	rel.dependencies[name] = ch
}

func (rel *Config) waitForDependencies() (err error) {
	for name, ch := range rel.dependencies {
		log.Debugf("release %s is waiting for dependency %s", rel.Name, name)
		status := <-ch
		log.Debugf("dependency %s of release %s done", name, rel.Name)
		if status == pubsub.ReleaseFailed {
			err = DependencyFailedError
		}
	}
	return
}

func (rel *Config) HandleDependencies() {
	for _, dep := range rel.DependsOn {
		rel.addDependency(dep)
	}
}
