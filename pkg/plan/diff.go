package plan

import (
	"errors"
	"os"

	"github.com/databus23/helm-diff/diff"
	"github.com/databus23/helm-diff/manifest"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
	live "helm.sh/helm/v3/pkg/release"
)

var ErrPlansAreTheSame = errors.New("plan1 and plan2 are the same")

// DiffPlan show diff between 2 plans
func (p *Plan) DiffPlan(b *Plan, showSecret bool, diffWide int) {
	visited := make([]uniqname.UniqName, 0, len(p.body.Releases))
	k := 0

	for _, rel := range p.body.Releases {
		m := rel.Uniq() + ".yml"
		visited = append(visited, rel.Uniq())

		oldSpecs := manifest.Parse(b.manifests[m], rel.Namespace)
		newSpecs := manifest.Parse(p.manifests[m], rel.Namespace)

		change := diff.Manifests(oldSpecs, newSpecs, []string{}, showSecret, diffWide, os.Stdout)
		if !change {
			k++
			log.Info("🆚 ❎ ", rel.Uniq(), " no changes")
		}
	}

	showChangesReport(p.body.Releases, visited, k)
}

// DiffLive show diff with production releases in k8s-cluster
func (p *Plan) DiffLive(showSecret bool, diffWide int) {
	alive, _, err := p.GetLive()
	if err != nil {
		log.Fatalf("Something went wrong with getting realeases in the kubernetes cluster: %v", err)
	}

	visited := make([]uniqname.UniqName, 0, len(p.body.Releases))
	k := 0
	for _, rel := range p.body.Releases {
		m := rel.Uniq() + ".yml"
		visited = append(visited, rel.Uniq())
		if active, ok := alive[rel.Uniq()]; ok {
			// I dont use manifest.ParseRelease
			// Because Structs are different.
			oldSpecs := manifest.Parse(active.Manifest, rel.Namespace)
			newSpecs := manifest.Parse(p.manifests[m], rel.Namespace)

			change := diff.Manifests(oldSpecs, newSpecs, []string{}, showSecret, diffWide, os.Stdout)
			if !change {
				k++
				log.Info("🆚 ❎ ", rel.Uniq(), " no changes")
			}
		}
	}

	showChangesReport(p.body.Releases, visited, k)
}

// showChangesReport help function for reporting helm-diff
func showChangesReport(releases []*release.Config, visited []uniqname.UniqName, k int) {
	previous := false
	for _, rel := range releases {
		if !uniqname.Contains(rel.Uniq(), visited) {
			previous = true
			log.Warn("🆚 ", rel.Uniq(), " was found in previous plan but not affected in new")
		}
	}

	if k == len(releases) && !previous {
		log.Info("🆚 🌝 Plan has no changes")
	}
}

func (p *Plan) GetLiveOf(name uniqname.UniqName) (*live.Release, error) {
	for _, rel := range p.body.Releases {
		if rel.Uniq() == name {
			return rel.Get()
		}
	}

	return nil, errors.New("release 404")
}

// GetLive returns maps of releases in a k8s-cluster
func (p *Plan) GetLive() (found map[uniqname.UniqName]*live.Release, notFound []uniqname.UniqName, err error) {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	found = make(map[uniqname.UniqName]*live.Release)

	for i := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel *release.Config) {
			defer wg.Done()
			r, err := rel.Get()
			if err != nil {
				log.Warnf("I cant get realease from k8s: %v", err)
				notFound = append(notFound, rel.Uniq())
			} else {
				found[rel.Uniq()] = r
			}
		}(wg, p.body.Releases[i])
	}

	if err := wg.Wait(); err != nil {
		return nil, nil, err
	}

	return found, notFound, nil
}
