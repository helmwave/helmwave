package release

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *config) Status() (*release.Release, error) {
	client := action.NewStatus(rel.Cfg())
	client.ShowDescription = true

	return client.Run(rel.Name())
}
