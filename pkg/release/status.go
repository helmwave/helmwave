package release

import (
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *Release) Status() (*release.Release, error) {
	client := action.NewStatus(rel.Cfg())
	client.ShowDescription = true

	r, err := client.Run(rel.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to get status of release %s: %w", rel.Uniq(), err)
	}

	return r, nil
}
