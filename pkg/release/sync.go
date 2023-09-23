package release

import (
	"context"
	"io/fs"

	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *config) Sync(ctx context.Context, baseFS fs.FS) (*release.Release, error) {
	ctx = helper.ContextWithReleaseUniq(ctx, rel.Uniq())

	// Run hooks
	if rel.dryRun {
		err := rel.Lifecycle.RunPreBuild(ctx)
		if err != nil {
			return nil, err
		}

		defer func() {
			err := rel.Lifecycle.RunPostBuild(ctx)
			if err != nil {
				rel.Logger().Errorf("got an error from postbuild hooks: %v", err)
			}
		}()
	} else {
		err := rel.Lifecycle.RunPreUp(ctx)
		if err != nil {
			return nil, err
		}

		defer func() {
			err := rel.Lifecycle.RunPostUp(ctx)
			if err != nil {
				rel.Logger().Errorf("got an error from postup hooks: %v", err)
			}
		}()
	}

	r, err := rel.upgrade(ctx, baseFS)

	if err == nil && !rel.dryRun && rel.ShowNotes {
		rel.Logger().Infof("🗒️ release notes:\n%s", r.Info.Notes)
	}

	return r, err
}

func (rel *config) SyncDryRun(ctx context.Context, baseFS fs.FS) (*release.Release, error) {
	old := rel.dryRun
	defer rel.DryRun(old)
	rel.DryRun(true)

	return rel.Sync(ctx, baseFS)
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
