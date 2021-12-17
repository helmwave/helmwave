package release

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *config) Uninstall() (*release.UninstallReleaseResponse, error) {
	client := action.NewUninstall(rel.Cfg())
	client.Timeout = rel.Timeout

	return client.Run(rel.Name())
}
