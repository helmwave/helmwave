package release

import (
	"fmt"

	release "helm.sh/helm/v4/pkg/release/v1"
)

func (rel *config) Status() (*release.Release, error) {
	client := rel.newStatus()

	r, err := client.Run(rel.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to get status of release %s: %w", rel.Uniq(), err)
	}

	return asRelease(r)
}
