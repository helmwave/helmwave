package release

import (
	"context"

	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *config) Sync(ctx context.Context, runHooks bool) (r *release.Release, err error) {
	ctx = helper.ContextWithReleaseUniq(ctx, rel.Uniq())

	// Run hooks
	if runHooks {
		preHook, postHook := rel.syncLifecycleHooks()

		err = preHook(ctx)
		if err != nil {
			return nil, err
		}

		defer func() {
			lifecycleErr := postHook(ctx)
			if lifecycleErr != nil && err == nil {
				err = lifecycleErr
			}
		}()
	}

	r, err = rel.upgrade(ctx)
	if err != nil {
		return nil, err
	}

	if rel.dryRun {
		return r, nil
	}

	if rel.Tests.Enabled {
		err = rel.test()
		if err != nil {
			rel.Logger().Errorf("helm tests failed")

			return nil, err
		}
	}

	if !rel.HideNotes {
		rel.Logger().Infof("🗒️ release notes:\n%s", r.Info.Notes)
	}

	return r, nil
}

type lifecycleHook func(context.Context) error

func (rel *config) syncLifecycleHooks() (pre, post lifecycleHook) {
	if rel.dryRun {
		pre = rel.LifecycleF.RunPreBuild
		post = rel.LifecycleF.RunPostBuild
	} else {
		pre = rel.LifecycleF.RunPreUp
		post = rel.LifecycleF.RunPostUp
	}

	return
}

func (rel *config) SyncDryRun(ctx context.Context, runHooks bool) (*release.Release, error) {
	old := rel.dryRun
	if !old {
		defer rel.DryRun(old)
		rel.DryRun(true)
	}

	return rel.Sync(ctx, runHooks)
}
