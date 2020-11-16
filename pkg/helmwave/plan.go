package helmwave

import (
	"fmt"
	"github.com/zhilyaev/helmwave/pkg/release"
	"github.com/zhilyaev/helmwave/pkg/repo"
	"github.com/zhilyaev/helmwave/pkg/yml"
)

func (c *Config) PlanRepos() {
	c.Plan.Body.Repositories = repo.Plan(c.Plan.Body.Releases, c.Yml.Body.Repositories)
	if c.Debug {
		fmt.Println("🛠 Planned repositories:")
		yml.Print(c.Plan.Body.Repositories)
	}
}

func (c *Config) PlanReleases() {
	c.Plan.Body.Releases = release.Plan(c.Tags.Value(), c.Yml.Body.Releases)
	if c.Debug {
		fmt.Println("🛠 Planned releases:")
		yml.Print(c.Plan.Body.Releases)
	}
}

func (c *Config) RenderValues() {
	for i, rel := range c.Plan.Body.Releases {
		rel.RenderValues(c.Debug)
		// Make easy
		c.Plan.Body.Releases[i].Values = rel.Values
	}
}
