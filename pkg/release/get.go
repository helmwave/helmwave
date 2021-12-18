package release

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *config) Get() (*release.Release, error) {
	// IDK wtf is going on
	rel.cfg = nil
	client := action.NewGet(rel.Cfg())
	return client.Run(rel.Name())
}

func (rel *config) GetValues() (map[string]interface{}, error) {
	client := action.NewGetValues(rel.Cfg())
	return client.Run(rel.Name())
}
