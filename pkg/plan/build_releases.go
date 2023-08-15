package plan

import (
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func buildReleases(tags []string, releases []release.Config, matchAll bool) ([]release.Config, error) {
	plan := make([]release.Config, 0)

	for _, r := range releases {
		if checkTagInclusion(tags, r.Tags(), matchAll) {
			var err error
			plan, err = addToPlan(plan, r, releases)
			if err != nil {
				log.WithError(err).Error("failed to build releases plan")

				return nil, err
			}
		}
	}

	return plan, nil
}

func addToPlan(
	plan release.Configs,
	rel release.Config,
	releases release.Configs,
) (release.Configs, error) {
	if r, contains := plan.Contains(rel); contains {
		if r != rel {
			return nil, release.DuplicateReleasesError{Uniq: rel.Uniq()}
		} else {
			return plan, nil
		}
	}

	newPlan := plan
	newPlan = append(newPlan, rel)

	deps := rel.DependsOn()
	newDeps := make([]*release.DependsOnReference, 0, len(deps))

	for _, dep := range deps {
		l := rel.Logger().WithField("dependency", dep.Uniq())
		l.Trace("searching for dependency")

		r, found := releases.ContainsUniq(dep.Uniq())
		if found {
			var err error
			newPlan, err = addToPlan(newPlan, r, releases)
			if err != nil {
				return nil, err
			}

			newDeps = append(newDeps, dep)

			continue
		}

		if dep.Optional {
			l.Warn("can't find dependency in available releases, skipping")
		} else {
			l.Error("can't find required dependency")

			return nil, release.ErrDepFailed
		}
	}

	rel.SetDependsOn(newDeps)

	return newPlan, nil
}

func releaseNames(a []release.Config) (n []string) {
	for _, r := range a {
		n = append(n, r.Uniq().String())
	}

	return n
}

// checkTagInclusion checks where any of release tags are included in target tags.
func checkTagInclusion(targetTags, releaseTags []string, matchAll bool) bool {
	if len(targetTags) == 0 {
		return true
	}

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
