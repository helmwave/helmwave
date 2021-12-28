package release

import (
	"fmt"

	log "github.com/sirupsen/logrus"
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
		return nil, fmt.Errorf("failed to merge values for release %s: %w", rel.Uniq(), err)
	}

	// Template
	if rel.dryRun {
		log.Debugf("📄 %q template manifest ", rel.Uniq())

		r, err := rel.newInstall().Run(ch, vals)
		if err != nil {
			return nil, fmt.Errorf("failed to dry-run install %s: %w", rel.Uniq(), err)
		}

		return r, nil
	}

	// Install
	if !rel.isInstalled() {
		log.Debugf("🧐 Release %q does not exist. Installing it now.", rel.Uniq())

		r, err := rel.newInstall().Run(ch, vals)
		if err != nil {
			return nil, fmt.Errorf("failed to install %s: %w", rel.Uniq(), err)
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
