package release

import (
	"fmt"

	release "helm.sh/helm/v4/pkg/release/v1"
)

func (rel *config) Get(version int) (*release.Release, error) {
	client := rel.newGet()
	client.Version = version

	r, err := client.Run(rel.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to get release %s: %w", rel.Uniq(), err)
	}

	return asRelease(r)
}

func (rel *config) GetValues() (map[string]any, error) {
	client := rel.newGetValues()

	r, err := client.Run(rel.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to get release values of %s: %w", rel.Uniq(), err)
	}

	return r, nil
}
