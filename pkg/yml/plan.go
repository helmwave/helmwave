package yml

import (
	log "github.com/sirupsen/logrus"
	"github.com/zhilyaev/helmwave/pkg/release"
	"github.com/zhilyaev/helmwave/pkg/repo"
)

func (c *Config) Save(file string, tags []string, dir string) error {
	err := c.Plan(tags, dir)
	if err != nil {
		return err
	}
	return Save(file, &c)
}

func (c *Config) Plan(tags []string, dir string) (err error) {
	c.PlanReleases(tags)

	if err = c.PlanReleasesValues(dir); err != nil {
		return err
	}

	c.PlanRepos()

	return nil
}

func (c *Config) PlanRepos() {
	c.Repositories = repo.Plan(c.Releases, c.Repositories)
	names := make([]string, 0)
	for _, v := range c.Repositories {
		names = append(names, v.Name)
	}
	log.WithField("repositories", names).Info("🛠 Yml -> 🗄 repositories")
}

func (c *Config) PlanReleases(tags []string) {
	c.Releases = release.Plan(tags, c.Releases)
	names := make([]string, 0)
	for _, v := range c.Releases {
		names = append(names, v.UniqName())
	}
	log.WithField("releases", names).Info("🛠 Yml -> 🛥 releases")
}

func (c *Config) PlanReleasesValues(dir string) error {
	return release.PlanValues(c.Releases, dir)
}
