package release

import (
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *Config) upgrade() (*release.Release, error) {
	client := rel.newUpgrade()

	ch, err := rel.GetChart()
	if err != nil {
		return nil, err
	}

	// Values
	valuesFiles := make([]string, 0, len(rel.Values))
	for i := range rel.Values {
		valuesFiles = append(valuesFiles, rel.Values[i].Get())
	}

	valOpts := &values.Options{ValueFiles: valuesFiles}
	vals, err := valOpts.MergeValues(getter.All(rel.Helm()))
	if err != nil {
		return nil, err
	}

	// Install
	if !rel.isInstalled() || rel.dryRun {
		if rel.dryRun {
			log.Debugf("📄 Templating manifest %q ", rel.Uniq())
		} else {
			log.Debugf("🧐 Release %q does not exist. Installing it now.", rel.Uniq())
		}

		return rel.newInstall().Run(ch, vals)
	}

	// Upgrade
	return client.Run(rel.Name, ch, vals)
}
