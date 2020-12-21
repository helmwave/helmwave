package helmwave

import (
	log "github.com/sirupsen/logrus"
	"github.com/zhilyaev/helmwave/pkg/release"
	"github.com/zhilyaev/helmwave/pkg/repo"
	"github.com/zhilyaev/helmwave/pkg/yml"
)

func (c *Config) PlanRepos() (err error) {
	c.Plan.Body.Repositories, err = repo.Plan(c.Plan.Body.Releases, c.Yml.Body.Repositories)
	if err != nil {
		return err
	}
	names := make([]string, 0)
	for _, v := range c.Plan.Body.Repositories {
		names = append(names, v.Name)
	}
	log.WithField("repositories", names).Info("🛠 -> 🗄")
	return nil
}
func (c *Config) PlanReleases() {
	c.Plan.Body.Releases = release.Plan(c.Tags.Value(), c.Yml.Body.Releases)
	names := make([]string, 0)
	for _, v := range c.Plan.Body.Releases {
		names = append(names, v.UniqName())
	}
	log.WithField("releases", names).Info("🛠 -> 🛥")
}
func (c *Config) PlanReleasesValues() error {
	return release.PlanValues(c.Plan.Body.Releases, c.PlanDir)
}

func (c *Config) Planfile() error {
	c.InitPlan()
	err := c.ReadHelmWaveYml()
	if err != nil {
		return err
	}

	// General
	c.Plan.Body.Project = c.Yml.Body.Project
	c.Plan.Body.Version = c.Yml.Body.Version

	// Releases
	c.PlanReleases()

	// Values
	if err := c.PlanReleasesValues(); err != nil {
		return err
	}

	// Repos
	if err := c.PlanRepos(); err != nil {
		return err
	}

	// Save Plan
	return yml.Save(c.Plan.File, c.Plan.Body)
}

func (c *Config) InitPlan() {
	if c.PlanDir[len(c.PlanDir)-1:] != "/" {
		c.PlanDir += "/"
	}
	c.Plan.File = c.PlanDir + "planfile"
	log.Info("🛠 Your planfile is ", c.Plan.File)
}
