package plan

import (
	"context"
	"slices"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) buildReleases(ctx context.Context, o BuildOptions) ([]release.Config, error) {
	plan := make([]release.Config, 0)

	planAdderFunction := addToPlanWithDependencies
	if !o.EnableDependencies {
		planAdderFunction = addToPlanWithoutDependencies
	}

	for _, r := range p.body.Releases {
		if !checkTagInclusion(o.Tags, r.Tags(), o.MatchAll) {
			continue
		}

		var err error
		plan, err = planAdderFunction(ctx, plan, r, p.body.Releases)
		if err != nil {
			log.WithError(err).Error("failed to build releases plan")

			return nil, err
		}
	}

	return plan, nil
}

func addToPlan(
	ctx context.Context,
	plan release.Configs,
	rel release.Config,
) (_ release.Configs, err error) {
	if r, contains := plan.Contains(rel); contains {
		if r != rel {
			return nil, release.NewDuplicateError(rel.Uniq())
		} else {
			return plan, nil
		}
	}

	// Run hooks
	lifeCycle := rel.LifeCycle()
	err = lifeCycle.RunPreBuild(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		lifecycleErr := lifeCycle.RunPostBuild(ctx)
		if lifecycleErr != nil {
			log.Errorf("got an error from postbuild hooks: %v", lifecycleErr)
			if err == nil {
				err = lifecycleErr
			}
		}
	}()

	return append(plan, rel), nil
}

func addToPlanWithoutDependencies(
	ctx context.Context,
	plan release.Configs,
	rel release.Config,
	_ release.Configs,
) (release.Configs, error) {
	plan, err := addToPlan(ctx, plan, rel)
	if err != nil {
		return nil, err
	}

	rel.SetDependsOn([]*release.DependsOnReference{})

	return plan, nil
}

func addToPlanWithDependencies(
	ctx context.Context,
	plan release.Configs,
	rel release.Config,
	releases release.Configs,
) (release.Configs, error) {
	newPlan, err := addToPlan(ctx, plan, rel)
	if err != nil {
		return nil, err
	}

	deps := rel.DependsOn()
	newDeps := make([]*release.DependsOnReference, 0, len(deps))

	for _, dep := range deps {
		l := rel.Logger().WithField("dependency", dep.Uniq())
		l.Trace("searching for dependency")

		r, found := releases.ContainsUniq(dep.Uniq())
		if found {
			var err error
			newPlan, err = addToPlanWithDependencies(ctx, newPlan, r, releases)
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

func releaseNames(a []release.Config) []string {
	return helper.SlicesMap(a, func(r release.Config) string {
		return r.Uniq().String()
	})
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
		contains := slices.Contains(releaseTags, t)
		if !contains {
			return false
		}
	}

	return true
}

func checkAnyTagInclusion(targetTags, releaseTags []string) bool {
	for _, t := range targetTags {
		contains := slices.Contains(releaseTags, t)
		if contains {
			return true
		}
	}

	return false
}
