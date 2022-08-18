package release

import (
	"context"

	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *config) Sync(ctx context.Context) (*release.Release, error) {
	return rel.upgrade(ctx)
}

func (rel *config) SyncDryRun(ctx context.Context) (*release.Release, error) {
	old := rel.dryRun
	defer rel.DryRun(old)
	rel.DryRun(true)

	return rel.Sync(ctx)
}

func (rel *config) Cfg() *action.Configuration {
	if rel.cfg == nil {
		var err error
		rel.cfg, err = helper.NewCfg(rel.Namespace())
		if err != nil {
			rel.Logger().Fatal(err)

			return nil
		}
	}

	return rel.cfg
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
