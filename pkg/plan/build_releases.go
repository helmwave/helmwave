package plan

import (
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
)

func buildReleases(tags []string, releases []*release.Config, matchAll bool) (plan []*release.Config) {
	if len(tags) == 0 {
		return releases
	}

	releasesMap := make(map[uniqname.UniqName]*release.Config)

	for _, r := range releases {
		releasesMap[r.Uniq()] = r
	}

	for _, r := range releases {
		if checkTagInclusion(tags, r.Tags, matchAll) {
			plan = addToPlan(plan, r, releasesMap)
		}
	}

	return plan
}

func addToPlan(plan []*release.Config, rel *release.Config,
	releases map[uniqname.UniqName]*release.Config) []*release.Config {
	if rel.In(plan) {
		return plan
	}

	r := append(plan, rel) // nolint:gocritic

	for _, depName := range rel.DependsOn {
		depUN := uniqname.UniqName(depName)

		if dep, ok := releases[depUN]; ok {
			r = addToPlan(r, dep, releases)
		} else {
			log.Warnf("cannot find dependency %s in available releases, skipping it", depName)
		}
	}

	return r
}

func releaseNames(a []*release.Config) (n []string) {
	for _, r := range a {
		n = append(n, string(r.Uniq()))
	}

	return n
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
