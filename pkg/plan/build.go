package plan

import (
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

type buildOptions int

const (
	Release buildOptions = iota
	Values
	Repositories
)

func (p *Plan) Build() {

}

func (p *Plan) buildReleases(tags []string, matchAll bool) *Plan {
	if len(tags) == 0 {
		return p
	}

	var result []*release.Config

	uniqMap := p.uniqReleaseNames()

	for _, r := range p.body.Releases {
		if checkTagInclusion(tags, r.Tags, matchAll) {

		}
	}

	return p
}

// checkTagInclusion checks where any of release tags are included in target tags.
func checkTagInclusion(targetTags, releaseTags []string, matchAll bool) bool {
	for _, t := range targetTags {
		contains := helper.Contains(t, releaseTags)
		if matchAll && !contains {
			return false
		}
		if !matchAll && contains {
			return true
		}
	}

	return matchAll
}

func (p *Plan) addToPlan(rel *release.Config, m map[string]*release.Config) *Plan {
	if rel.In(p.body.Releases) {
		return p
	}

	p.body.Releases = append(p.body.Releases, rel)

	for _, depName := range rel.DependsOn {
		if dep, ok := m[depName]; ok {
			p.addToPlan(dep, m)
		} else {
			log.Warnf("cannot find dependency %s in available releases, skipping it", depName)
		}
	}

	return p
}

// uniqReleaseNames
func (p *Plan) uniqReleaseNames() map[string]*release.Config {
	m := make(map[string]*release.Config)
	for _, r := range p.body.Releases {
		m[r.UniqName()] = r
	}
	return m
}

// buildValues to planfile
func (p *Plan) buildValues() error {
	for _, rel := range p.body.Releases {
		err := rel.RenderValues(p.dir)
		if err != nil {
			return err
		}

		log.WithFields(log.Fields{
			"release":   rel.Name,
			"namespace": rel.Options.Namespace,
			"values":    rel.Values,
		}).Debug("🐞 Render Values")
	}

	return nil
}

// buildRepositories to planfile
func (p *Plan) buildRepositories() *Plan {
	return p
}
