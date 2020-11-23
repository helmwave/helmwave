package helmwave

import (
	log "github.com/sirupsen/logrus"
	"github.com/zhilyaev/helmwave/pkg/release"
	"github.com/zhilyaev/helmwave/pkg/repo"
)

func (c *Config) PlanRepos() {
	c.Plan.Body.Repositories = repo.Plan(c.Plan.Body.Releases, c.Yml.Body.Repositories)
}

func (c *Config) PlanReleases() {
	c.Plan.Body.Releases = release.Plan(c.Tags.Value(), c.Yml.Body.Releases)
}

func (c *Config) RenderValues() {
	for i, rel := range c.Plan.Body.Releases {
		rel.RenderValues()
		c.Plan.Body.Releases[i].Values = rel.Values
		log.Debugf("🐞 Values of %s : %+v", rel.Name, c.Plan.Body.Releases[i].Values)
	}
}
