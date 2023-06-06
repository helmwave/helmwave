package plan

import (
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
)

func buildReleases(tags []string, releases []release.Config, matchAll bool) ([]release.Config, error) {
	if len(tags) == 0 {
		return releases, nil
	}

	releasesMap := make(map[uniqname.UniqName]release.Config)

	for _, r := range releases {
		releasesMap[r.Uniq()] = r
	}

	plan := make([]release.Config, 0)

	for _, r := range releases {
		if checkTagInclusion(tags, r.Tags(), matchAll) {
			var err error
			plan, err = addToPlan(plan, r, releasesMap)
			if err != nil {
				log.WithError(err).Error("failed to build releases plan")

				return nil, err
			}
		}
	}

	return plan, nil
}

//nolintlint:nestif
func addToPlan(plan []release.Config, rel release.Config,
	releases map[uniqname.UniqName]release.Config,
) ([]release.Config, error) {
	if helper.In(rel, plan) {
		return plan, nil
	}

	r := plan
	r = append(r, rel)

	for _, dep := range rel.DependsOn() {
		if depRel, ok := releases[dep.Uniq()]; ok {
			var err error
			r, err = addToPlan(r, depRel, releases)
			if err != nil {
				return nil, err
			}
		} else {
			if dep.Optional {
				log.Warnf("cannot find dependency %q in available releases, skipping it", dep.Uniq())
			} else {
				rel.Logger().WithField("dependency", dep.Uniq()).Error("cannot find required dependency")

				return nil, release.ErrDepFailed
			}
		}
	}

	return r, nil
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
