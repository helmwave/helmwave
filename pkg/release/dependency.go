package release

import (
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"sort"
	"time"

	"github.com/helmwave/helmwave/pkg/pubsub"
	log "github.com/sirupsen/logrus"
)

var releasePubSub = pubsub.NewReleasePubSub()

func (rel *Config) NotifySuccess() {
	if !rel.dryRun {
		releasePubSub.PublishSuccess(rel.UniqName())
	}
}

func (rel *Config) NotifyFailed() {
	if !rel.dryRun {
		releasePubSub.PublishFailed(rel.UniqName())
	}
}

func (rel *Config) addDependency(name uniqname.UniqName) {
	ch := releasePubSub.Subscribe(name)

	if rel.dependencies == nil {
		rel.dependencies = make(map[uniqname.UniqName]<-chan pubsub.ReleaseStatus)
	}

	rel.dependencies[name] = ch
}

func (rel *Config) waitForDependencies() (err error) {
	if rel.dryRun {
		return nil
	}

	for name, ch := range rel.dependencies {
		status := rel.waitForDependency(ch, name)
		if status == pubsub.ReleaseFailed {
			err = ErrDepFailed
		}
	}
	return
}

func (rel *Config) waitForDependency(ch <-chan pubsub.ReleaseStatus, name uniqname.UniqName) pubsub.ReleaseStatus {
	ticker := time.NewTicker(5 * time.Second)
	var status pubsub.ReleaseStatus

F:
	for {
		select {
		case status = <-ch:
			ticker.Stop()
			break F
		case <-ticker.C:
			log.Infof("release %s is waiting for dependency %s", rel.UniqName(), name)
		}
	}
	log.Infof("dependency %s of release %s done", name, rel.UniqName())
	return status
}

func (rel *Config) HandleDependencies(releases []*Config) {
	sort.Strings(rel.DependsOn)

	depsAdded := make(map[string]bool)
	for _, r := range releases {
		name := r.UniqName()
		if i := sort.SearchStrings(rel.DependsOn, string(name)); i < len(rel.DependsOn) && rel.DependsOn[i] == string(name) {
			rel.addDependency(name)
			depsAdded[string(name)] = true
		}
	}

	for _, dep := range rel.DependsOn {
		if !depsAdded[dep] {
			log.Warnf("cannot find dependency %s in plan, skipping it", dep)
		}
	}
}
