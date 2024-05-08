package release

import (
	"fmt"

	"helm.sh/helm/v3/pkg/release"
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
		return result[0], nil
	default:
		return nil, ErrFoundMultiple
	}
}
