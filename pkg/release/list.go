package release

import (
	"fmt"

	release "helm.sh/helm/v4/pkg/release/v1"
)

func (rel *config) List() (*release.Release, error) {
	client := rel.newList()

	result, err := client.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to list release %s: %w", rel.Uniq(), err)
	}

	switch len(result) {
	case 0:
		return nil, ErrNotFound
	case 1:
		return asRelease(result[0])
	default:
		return nil, ErrFoundMultiple
	}
}
