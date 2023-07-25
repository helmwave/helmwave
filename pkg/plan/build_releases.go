package plan

import (
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
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

//nolintlint:nestif
//nolint:gocognit
func addToPlan(plan []release.Config, rel release.Config,
	releases []release.Config,
) ([]release.Config, error) {
	for _, r := range plan {
		if r.Uniq() == rel.Uniq() {
			if r != rel {
				return nil, DuplicateReleasesError{uniq: rel.Uniq()}
			} else {
				return plan, nil
			}
		}
	}

	r := plan
	r = append(r, rel)

	deps := rel.DependsOn()

	for i, dep := range deps {
		found := false
		l := rel.Logger().WithField("dependency", dep.Uniq())
		l.Trace("searching for dependency")

		for _, rel := range releases {
			if rel.Uniq().Equal(dep.Uniq()) {
				found = true
				var err error
				r, err = addToPlan(r, rel, releases)
				if err != nil {
					return nil, err
				}
			}
		}

		if !found {
			if dep.Optional {
				l.Warn("can't find dependency in available releases, skipping")
				deps = slices.Delete(deps, i, i+1)
			} else {
				l.Error("can't find required dependency")

				return nil, release.ErrDepFailed
			}
		}
	}

	rel.SetDependsOn(deps)

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
