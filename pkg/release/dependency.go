package release

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/zhilyaev/helmwave/pkg/pubsub"
	"sort"
)

var releasePubSub = pubsub.NewReleasePubSub()
var DependencyFailedError = errors.New("dependency failed")

func (rel *Config) NotifySuccess() {
	if !rel.Options.DryRun {
		releasePubSub.PublishSuccess(rel.Name)
	}
}

func (rel *Config) NotifyFailed() {
	if !rel.Options.DryRun {
		releasePubSub.PublishFailed(rel.Name)
	}
}

func (rel *Config) addDependency(name string) {
	ch := releasePubSub.Subscribe(name)

	if rel.dependencies == nil {
		rel.dependencies = make(map[string]<-chan pubsub.ReleaseStatus)
	}

	rel.dependencies[name] = ch
}

func (rel *Config) waitForDependencies() (err error) {
	if rel.Options.DryRun {
		return nil
	}

	for name, ch := range rel.dependencies {
		log.Infof("release %s is waiting for dependency %s", rel.Name, name)
		status := <-ch
		log.Infof("dependency %s of release %s done", name, rel.Name)
		if status == pubsub.ReleaseFailed {
			err = DependencyFailedError
		}
	}
	return
}

func (rel *Config) HandleDependencies(releases []*Config) {
	sort.Strings(rel.DependsOn)

	for _, r := range releases {
		if i := sort.SearchStrings(rel.DependsOn, r.Name); i < len(rel.DependsOn) && rel.DependsOn[i] == r.Name {
			rel.addDependency(r.Name)
		}
	}
}
