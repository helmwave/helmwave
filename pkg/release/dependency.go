package release

import (
	"errors"
	"sort"
	"time"

	"github.com/helmwave/helmwave/pkg/pubsub"
	log "github.com/sirupsen/logrus"
)

var (
	releasePubSub         = pubsub.NewReleasePubSub()
	DependencyFailedError = errors.New("dependency failed")
)

func (rel *Config) NotifySuccess() {
	if !rel.Options.DryRun {
		releasePubSub.PublishSuccess(rel.UniqName())
	}
}

func (rel *Config) NotifyFailed() {
	if !rel.Options.DryRun {
		releasePubSub.PublishFailed(rel.UniqName())
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
		status := rel.waitForDependency(ch, name)
		if status == pubsub.ReleaseFailed {
			err = DependencyFailedError
		}
	}
	return
}

func (rel *Config) waitForDependency(ch <-chan pubsub.ReleaseStatus, name string) pubsub.ReleaseStatus {
	ticker := time.NewTicker(5 * time.Second)
	var status pubsub.ReleaseStatus

F:
	for {
		select {
		case status = <-ch:
			ticker.Stop()
			break F
		case <-ticker.C:
			log.Infof("release %s is waiting for dependency %s", rel.Name, name)
		}
	}
	log.Infof("dependency %s of release %s done", name, rel.Name)
	return status
}

func (rel *Config) HandleDependencies(releases []*Config) {
	sort.Strings(rel.DependsOn)

	depsAdded := make(map[string]bool)
	for _, r := range releases {
		name := r.UniqName()
		if i := sort.SearchStrings(rel.DependsOn, name); i < len(rel.DependsOn) && rel.DependsOn[i] == name {
			rel.addDependency(name)
			depsAdded[name] = true
		}
	}

	for _, dep := range rel.DependsOn {
		if !depsAdded[dep] {
			log.Warnf("cannot find dependency %s in plan, skipping it", dep)
		}
	}
}
