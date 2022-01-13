package release

import (
	"fmt"

	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *config) upgrade() (*release.Release, error) {
	client := rel.newUpgrade()

	ch, err := rel.GetChart()
	if err != nil {
		return nil, err
	}

	// Values
	valuesFiles := make([]string, 0, len(rel.Values()))
	for i := range rel.Values() {
		valuesFiles = append(valuesFiles, rel.Values()[i].Get())
	}

	valOpts := &values.Options{ValueFiles: valuesFiles}
	vals, err := valOpts.MergeValues(getter.All(rel.Helm()))
	if err != nil {
		return nil, fmt.Errorf("failed to merge values for release %q: %w", rel.Uniq(), err)
	}

	// Template
	if rel.dryRun {
		rel.Logger().Debug("📄 template manifest")

		r, err := rel.newInstall().Run(ch, vals)
		if err != nil {
			return nil, fmt.Errorf("failed to dry-run install %q: %w", rel.Uniq(), err)
		}

		return r, nil
	}

	// Install
	if !rel.isInstalled() {
		rel.Logger().Debug("🧐 Release does not exist. Installing it now.")

		r, err := rel.newInstall().Run(ch, vals)
		if err != nil {
			return nil, fmt.Errorf("failed to install %q: %w", rel.Uniq(), err)
		}

		return r, nil
	}

	// Upgrade
	r, err := client.Run(rel.Name(), ch, vals)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade %s: %w", rel.Uniq(), err)
	}

	return r, nil
}
