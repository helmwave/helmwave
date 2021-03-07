package release

import (
	"errors"
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
	rel.dependencies = append(rel.dependencies, ch)
}

func (rel *Config) waitForDependencies() (err error) {
	for _, ch := range rel.dependencies {
		status := <-ch
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
