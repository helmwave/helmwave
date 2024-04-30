package release

import (
	"fmt"

	"helm.sh/helm/v3/pkg/release"
)

func (rel *config) Status() (*release.Release, error) {
	client := rel.newStatus()

	r, err := client.Run(rel.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to get status of release %s: %w", rel.Uniq(), err)
	}

	return r, nil
}
