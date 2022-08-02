package release

import (
	"context"
	"fmt"

	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *config) upgrade(ctx context.Context) (*release.Release, error) {
	client := rel.newUpgrade()

	ch, err := rel.GetChart()
	if err != nil {
		return nil, err
	}

	// Values
	valuesFiles := make([]string, 0, len(rel.Values()))
	for i := range rel.Values() {
		valuesFiles = append(valuesFiles, rel.Values()[i].Dst)
	}

	valOpts := &values.Options{ValueFiles: valuesFiles}
	vals, err := valOpts.MergeValues(getter.All(rel.Helm()))
	if err != nil {
		return nil, fmt.Errorf("failed to merge values %v: %w", valuesFiles, err)
	}

	// Install
	if !rel.isInstalled() {
		if !rel.dryRun {
			rel.Logger().Debug("🧐 Release does not exist. Installing it now.")
		}

		r, err := rel.newInstall().RunWithContext(ctx, ch, vals)
		if err != nil {
			return nil, fmt.Errorf("failed to install %q: %w", rel.Uniq(), err)
		}

		return r, nil
	}

	// Upgrade
	r, err := client.RunWithContext(ctx, rel.Name(), ch, vals)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade %s: %w", rel.Uniq(), err)
	}

	return r, nil
}
