package release

import (
	"sort"
	"time"

	"github.com/helmwave/helmwave/pkg/pubsub"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
)

// TODO: we need to move this out of global context.
var releasePubSub = pubsub.NewReleasePubSub()

func (rel *config) NotifySuccess() {
	if rel.dryRun {
		return
	}
	releasePubSub.PublishSuccess(rel.Uniq())
}

func (rel *config) NotifyFailed() {
	if rel.dryRun {
		return
	}

	if rel.AllowFailure {
		rel.Logger().Warn("failed but is allowed to fail")
		releasePubSub.PublishSuccess(rel.Uniq())

		return
	}

	releasePubSub.PublishFailed(rel.Uniq())
}

func (rel *config) addDependency(name uniqname.UniqName) {
	ch := releasePubSub.Subscribe(name)

	if rel.dependencies == nil {
		rel.dependencies = make(map[uniqname.UniqName]<-chan pubsub.ReleaseStatus)
	}

	rel.dependencies[name] = ch
}

func (rel *config) waitForDependencies() (err error) {
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

func (rel *config) waitForDependency(ch <-chan pubsub.ReleaseStatus, name uniqname.UniqName) pubsub.ReleaseStatus {
	ticker := time.NewTicker(5 * time.Second)
	var status pubsub.ReleaseStatus

F:
	for {
		select {
		case status = <-ch:
			ticker.Stop()

			break F
		case <-ticker.C:
			rel.Logger().Infof("waiting for dependency %s", name)
		}
	}
	rel.Logger().Infof("dependency %s done", name)

	return status
}

func (rel *config) HandleDependencies(releases []Config) {
	sort.Strings(rel.DependsOn())

	depsAdded := make(map[string]bool)
	for _, r := range releases {
		name := r.Uniq()
		i := sort.SearchStrings(rel.DependsOn(), string(name))
		if i < len(rel.DependsOn()) && rel.DependsOn()[i] == string(name) {
			rel.addDependency(name)
			depsAdded[string(name)] = true
		}
	}

	for _, dep := range rel.DependsOn() {
		if !depsAdded[dep] {
			rel.Logger().Warnf("cannot find dependency %s in plan, skipping it", dep)
		}
	}
}
