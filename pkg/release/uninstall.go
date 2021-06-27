package release

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *Config) uninstall() (*release.UninstallReleaseResponse, error) {
	cfg, err := rel.cfg()
	if err != nil {
		return nil, err
	}

	client := action.NewUninstall(cfg)
	return client.Run(rel.ReleaseName)
}
