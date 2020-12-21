package yml

import (
	log "github.com/sirupsen/logrus"
	"github.com/zhilyaev/helmwave/pkg/release"
	"github.com/zhilyaev/helmwave/pkg/repo"
)

func (c *Config) Plan(tags []string, dir string) error {
	err := c.Body.Plan(tags, dir)
	if err != nil {
		return err
	}
	return c.Save()
}

func (c *Config) Save() error {
	return Save(c.File, c.Body)
}

func (c *Body) Plan(tags []string, dir string) (err error) {
	c.PlanReleases(tags)

	if err = c.PlanReleasesValues(dir); err != nil {
		return err
	}

	return c.PlanRepos()
}

func (c *Body) PlanRepos() (err error) {
	c.Repositories = repo.Plan(c.Releases, c.Repositories)
	names := make([]string, 0)
	for _, v := range c.Repositories {
		names = append(names, v.Name)
	}
	log.WithField("repositories", names).Info("🛠 Plan -> 🗄 repositories")
	return nil
}

func (c *Body) PlanReleases(tags []string) {
	c.Releases = release.Plan(tags, c.Releases)
	names := make([]string, 0)
	for _, v := range c.Releases {
		names = append(names, v.UniqName())
	}
	log.WithField("releases", names).Info("🛠 Plan -> 🛥 releases")
}

func (c *Body) PlanReleasesValues(dir string) error {
	return release.PlanValues(c.Releases, dir)
}
