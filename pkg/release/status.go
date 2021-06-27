package release

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *Config) Status() (*release.Release, error) {
	cfg, err := rel.cfg()
	if err != nil {
		return nil, err
	}

	client := action.NewStatus(cfg)
	client.ShowDescription = true

	return client.Run(rel.ReleaseName)
}
