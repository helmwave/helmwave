package release

import (
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *config) Get() (*release.Release, error) {
	client := action.NewGet(rel.Cfg())

	r, err := client.Run(rel.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to get release %s: %w", rel.Uniq(), err)
	}

	return r, nil
}

func (rel *config) GetValues() (map[string]any, error) {
	client := action.NewGetValues(rel.Cfg())

	r, err := client.Run(rel.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to get release values of %s: %w", rel.Uniq(), err)
	}

	return r, nil
}
