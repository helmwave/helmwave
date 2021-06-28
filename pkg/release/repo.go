package release

import (
	"helm.sh/helm/v3/pkg/chart/loader"
	"strings"
)

// Repositories returns repository for release
func (rel *Config) Repositories() (repos []string, err error) {
	chart, err := loader.Load(rel.Chart.Name)
	if err != nil {
		return nil, err
	}

	deps := chart.Metadata.Dependencies

	for _, d := range deps {
		d.Repository = strings.TrimPrefix(d.Repository, "@")
		repos = append(repos, d.Repository)
	}

	return repos, nil
}
