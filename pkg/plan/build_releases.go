package plan

import (
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func buildReleases(tags []string, releases []release.Config, matchAll bool) ([]release.Config, error) {
	if len(tags) == 0 {
		return releases, nil
	}

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

	for _, dep := range rel.DependsOn() {
		found := false
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
				log.Warnf("can't find dependency %q in available releases, skipping", dep.Uniq())
			} else {
				rel.Logger().WithField("dependency", dep.Uniq()).Error("can't find required dependency")

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
