package release

import (
	"errors"
	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"strings"
)

// RepoDeps returns repository for release
func (rel *Config) RepoDeps() (repos []string, err error) {
	err = rel.ChartDepsUpd()
	if err != nil {
		log.Warn("Cant get deps for ", rel.Uniq())
	}

	chart, err := rel.GetChart()
	if err != nil {
		log.Warn("Failed get chart")
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
