package release

import (
	"context"

	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

//nolint:gocognit,nestif,cyclop
func (rel *config) Sync(ctx context.Context, runHooks bool) (r *release.Release, err error) {
	ctx = helper.ContextWithReleaseUniq(ctx, rel.Uniq())

	// Run hooks
	if runHooks {
		if rel.dryRun {
			err = rel.Lifecycle.RunPreBuild(ctx)
			if err != nil {
				return
			}

			defer func() {
				lifecycleErr := rel.Lifecycle.RunPostBuild(ctx)
				if lifecycleErr != nil {
					rel.Logger().Errorf("got an error from postbuild hooks: %v", lifecycleErr)
					if err == nil {
						err = lifecycleErr
					}
				}
			}()
		} else {
			err = rel.Lifecycle.RunPreUp(ctx)
			if err != nil {
				return
			}

			defer func() {
				lifecycleErr := rel.Lifecycle.RunPostUp(ctx)
				if lifecycleErr != nil {
					rel.Logger().Errorf("got an error from postup hooks: %v", lifecycleErr)
					if err == nil {
						err = lifecycleErr
					}
				}
			}()
		}
	}

	r, err = rel.upgrade(ctx)
	if err != nil {
		return
	}

	if rel.Tests.Enabled && !rel.dryRun {
		err = rel.test()
		if err != nil {
			rel.Logger().Errorf("helm tests failed")

			return
		}
	}

	if !rel.dryRun && rel.ShowNotes {
		rel.Logger().Infof("🗒️ release notes:\n%s", r.Info.Notes)
	}

	return
}

func (rel *config) SyncDryRun(ctx context.Context, runHooks bool) (*release.Release, error) {
	old := rel.dryRun
	defer rel.DryRun(old)
	rel.DryRun(true)

	return rel.Sync(ctx, runHooks)
}

func (rel *config) Cfg() *action.Configuration {
	cfg, err := helper.NewCfg(rel.Namespace(), rel.KubeContext())
	if err != nil {
		rel.Logger().Fatal(err)

		return nil
	}

	return cfg
}

func (rel *config) Helm() *helm.EnvSettings {
	if rel.helm == nil {
		var err error
		rel.helm, err = helper.NewHelm(rel.Namespace())
		if err != nil {
			rel.Logger().Fatal(err)

			return nil
		}

		rel.helm.Debug = helper.Helm.Debug
	}

	return rel.helm
}
