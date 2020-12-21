package helmwave

import (
	log "github.com/sirupsen/logrus"
	"github.com/zhilyaev/helmwave/pkg/release"
	"github.com/zhilyaev/helmwave/pkg/repo"
	"github.com/zhilyaev/helmwave/pkg/yml"
)

func (c *Config) PlanRepos() (err error) {
	c.Plan.Body.Repositories, err = repo.Plan(c.Plan.Body.Releases, c.Yml.Body.Repositories)
	return err
}

func (c *Config) PlanReleases() {
	c.Plan.Body.Releases = release.Plan(c.Tags.Value(), c.Yml.Body.Releases)
}

func (c *Config) PlanValues() error {
	for i, rel := range c.Plan.Body.Releases {
		err := rel.RenderValues(c.PlanDir)
		if err != nil {
			return err
		}

		c.Plan.Body.Releases[i].Values = rel.Values
		log.WithFields(log.Fields{
			"release":   rel.Name,
			"namespace": rel.Options.Namespace,
			"values":    c.Plan.Body.Releases[i].Values,
		}).Debug("🐞 Render Values")
	}

	return nil
}

func (c *Config) InitPlanDirFile() {
	if c.PlanDir[len(c.PlanDir)-1:] != "/" {
		c.PlanDir += "/"
	}
	c.Plan.File = c.PlanDir + "planfile"
	log.Info("🛠 Your planfile is ", c.Plan.File)
}

func (c *Config) GenPlanfile() error {
	c.InitPlanDirFile()
	err := c.ReadHelmWaveYml()
	if err != nil {
		return err
	}

	c.Plan.Body.Project = c.Yml.Body.Project
	c.Plan.Body.Version = c.Yml.Body.Version

	// Releases
	c.PlanReleases()
	if err := c.PlanValues(); err != nil {
		return err
	}

	names := make([]string, 0)
	for _, v := range c.Plan.Body.Releases {
		names = append(names, v.Name)
	}
	log.WithField("releases", names).Info("🛠 -> 🛥")

	// Repos
	c.PlanRepos()
	names = make([]string, 0)
	for _, v := range c.Plan.Body.Repositories {
		names = append(names, v.Name)
	}
	log.WithField("repositories", names).Info("🛠 -> 🗄")

	return yml.Save(c.Plan.File, c.Plan.Body)
}
