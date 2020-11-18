package helmwave

import (
	"fmt"
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
		rel.RenderValues(c.Debug)
		c.Plan.Body.Releases[i].Values = rel.Values
		if c.Debug {
			fmt.Println("🐞 Values of", rel.Name, ":", c.Plan.Body.Releases[i].Values)
		}
	}
}
