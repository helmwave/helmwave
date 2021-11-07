package release

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *Config) Get() (*release.Release, error) {
	client := action.NewGet(rel.Cfg())
	return client.Run(rel.Name)
}

func (rel *Config) GetValues() (map[string]interface{}, error) {
	client := action.NewGetValues(rel.Cfg())
	return client.Run(rel.Name)
}
