package release

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *Config) Status() (*release.Release, error) {
	var err error
	rel.cfg, err = rel.newCfg()
	if err != nil {
		return nil, err
	}

	client := action.NewStatus(rel.cfg)
	client.ShowDescription = true

	return client.Run(rel.Name)
}
