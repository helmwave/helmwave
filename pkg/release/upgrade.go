package release

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
)

// Helm wraps a lot of meta.NoKindMatchError into fmt.Errorf which makes errors.Is unusable.
// So we have to find this substring in error string.
const errMissingCRD = "unable to build kubernetes objects from release manifest:"

func (rel *config) upgrade(ctx context.Context, baseFS fs.FS) (*release.Release, error) {
	ch, err := rel.GetChart(baseFS)
	if err != nil {
		return nil, err
	}

	// Values
	valuesFiles := make([]string, 0, len(rel.Values()))
	for i := range rel.Values() {
		valuesFiles = append(valuesFiles, rel.Values()[i].Src)
	}

	valOpts := &values.Options{ValueFiles: valuesFiles}
	vals, err := valOpts.MergeValues(getter.All(rel.Helm()))
	if err != nil {
		return nil, fmt.Errorf("failed to merge values %v: %w", valuesFiles, err)
	}

	// Install or Template
	if rel.dryRun {
		rel.Logger().Debug("I'll dry-run.")
		r, err := rel.installWithRetry(ctx, ch, vals)
		if err != nil {
			return nil, fmt.Errorf("failed with dry-run %q: %w", rel.Uniq(), err)
		}

		return r, nil
	} else if !rel.dryRun && !rel.isInstalled() {
		rel.Logger().Debug("🧐 Release does not exist. Installing it now.")
		r, err := rel.installWithRetry(ctx, ch, vals)
		if err != nil {
			return nil, fmt.Errorf("failed to install %q: %w", rel.Uniq(), err)
		}

		return r, nil
	}

	pending, err := rel.isPending()
	if err != nil {
		return nil, fmt.Errorf("failed to check %q for pending status: %w", rel.Uniq(), err)
	}
	if pending {
		err := rel.fixPending(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fix %q pending status: %w", rel.Uniq(), err)
		}
	}

	// Upgrade
	r, err := rel.upgradeWithRetry(ctx, ch, vals)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade %s: %w", rel.Uniq(), err)
	}

	return r, nil
}

//nolint:wrapcheck // we wrap it later
func (rel *config) installWithRetry(
	ctx context.Context,
	ch *chart.Chart,
	vals map[string]interface{},
) (*release.Release, error) {
	r, err := rel.newInstall().RunWithContext(ctx, ch, vals)

	if err != nil && strings.Contains(err.Error(), errMissingCRD) && rel.dryRun {
		er := rel.forceOfflineKubeVersion()
		// return original error if we can't get kubernetes version
		if er != nil {
			return r, err
		}

		return rel.newInstall().RunWithContext(ctx, ch, vals)
	}

	return r, err
}

//nolint:wrapcheck // we wrap it later
func (rel *config) upgradeWithRetry(
	ctx context.Context,
	ch *chart.Chart,
	vals map[string]interface{},
) (*release.Release, error) {
	r, err := rel.newUpgrade().RunWithContext(ctx, rel.Name(), ch, vals)

	if err != nil && strings.Contains(err.Error(), errMissingCRD) && rel.dryRun {
		er := rel.forceOfflineKubeVersion()
		// return original error if we can't get kubernetes version
		if er != nil {
			return r, err
		}

		return rel.newUpgrade().RunWithContext(ctx, rel.Name(), ch, vals)
	}

	return r, err
}

func (rel *config) forceOfflineKubeVersion() error {
	rel.Logger().Warn("🤔hmm, it looks like some required CRDs are not installed, setting offline_kube_version and trying again")

	v, err := helper.GetKubernetesVersion(rel.Cfg())
	if err != nil {
		rel.Logger().WithError(err).Error("cannot get current kubernetes version, you need to set it manually")

		return err
	}

	rel.OfflineKubeVersionF = v.GitVersion
	rel.Logger().WithField("version", rel.OfflineKubeVersionF).Info("discovered kubernetes version")

	return nil
}
