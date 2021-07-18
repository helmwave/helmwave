package release

import (
	"errors"
	"github.com/helmwave/helmwave/pkg/helper"
	"strings"
)

// RepoDeps returns repository for release
func (rel *Config) RepoDeps() (repos []string, err error) {

	chart, err := rel.GetChart()
	if err != nil {
		return nil, err
	}

	repos = append(repos, rel.Repo())

	for _, d := range chart.Metadata.Dependencies {
		if d.Enabled {
			d.Repository = strings.TrimPrefix(d.Repository, "@")
			repos = append(repos, d.Repository)
		} else if helper.IsURL(d.Repository) {
			return nil, errors.New("url is unsupported: " + d.Repository)
		}
	}

	return repos, nil
}

func (rel *Config) Repo() string {
	return strings.Split(rel.Chart.Name, "/")[0]
}
