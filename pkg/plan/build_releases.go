package plan

import (
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
)

func buildReleases(tags []string, releases []release.Config, matchAll bool) (plan []release.Config) {
	if len(tags) == 0 {
		return releases
	}

	releasesMap := make(map[uniqname.UniqName]release.Config)

	for _, r := range releases {
		releasesMap[r.Uniq()] = r
	}

	for _, r := range releases {
		if checkTagInclusion(tags, r.Tags(), matchAll) {
			plan = addToPlan(plan, r, releasesMap)
		}
	}

	return plan
}

func addToPlan(plan []release.Config, rel release.Config,
	releases map[uniqname.UniqName]release.Config,
) []release.Config {
	if helper.In(rel, plan) {
		return plan
	}

	r := plan
	r = append(r, rel)

	for _, depName := range rel.DependsOn() {
		if dep, ok := releases[depName]; ok {
			r = addToPlan(r, dep, releases)
		} else {
			log.Warnf("cannot find dependency %q in available releases, skipping it", depName)
		}
	}

	return r
}

func releaseNames(a []release.Config) (n []string) {
	for _, r := range a {
		n = append(n, r.Uniq().String())
	}

	return n
}

// checkTagInclusion checks where any of release tags are included in target tags.
func checkTagInclusion(targetTags, releaseTags []string, matchAll bool) bool {
	if matchAll {
		return checkAllTagsInclusion(targetTags, releaseTags)
	}

	return checkAnyTagInclusion(targetTags, releaseTags)
}

func checkAllTagsInclusion(targetTags, releaseTags []string) bool {
	for _, t := range targetTags {
		contains := helper.Contains(t, releaseTags)
		if !contains {
			return false
		}
	}

	return true
}

func checkAnyTagInclusion(targetTags, releaseTags []string) bool {
	for _, t := range targetTags {
		contains := helper.Contains(t, releaseTags)
		if contains {
			return true
		}
	}

	return false
}
