package release

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *Config) uninstall() (*release.UninstallReleaseResponse, error) {
	client := action.NewUninstall(rel.cfg)
	client.Timeout = rel.Timeout
	client.DryRun = rel.dryRun

	return client.Run(rel.Name)
}
