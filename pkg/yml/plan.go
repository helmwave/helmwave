package yml

import (
	"github.com/helmwave/helmwave/pkg/feature"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	log "github.com/sirupsen/logrus"
	helm "helm.sh/helm/v3/pkg/cli"
)

type SavePlanOptions struct {
	file string
	tags []string
	dir  string

	planReleases bool
	planRepos    bool
	planValues   bool
}

func (o *SavePlanOptions) File(file string) *SavePlanOptions {
	o.file = file
	return o
}

func (o *SavePlanOptions) Tags(tags []string) *SavePlanOptions {
	o.tags = tags
	return o
}

func (o *SavePlanOptions) Dir(dir string) *SavePlanOptions {
	o.dir = dir
	return o
}

func (o *SavePlanOptions) PlanReleases() *SavePlanOptions {
	o.planReleases = true
	return o.PlanValues().PlanRepos()
}

func (o *SavePlanOptions) PlanRepos() *SavePlanOptions {
	o.planRepos = true
	return o
}

func (o *SavePlanOptions) PlanValues() *SavePlanOptions {
	o.planValues = true
	return o
}

func (c *Config) SavePlan(o *SavePlanOptions, helmSettings *helm.EnvSettings) error {
	err := c.Plan(o, helmSettings)
	if err != nil {
		return err
	}
	return Save(o.file, &c)
}

func (c *Config) Plan(o *SavePlanOptions, helmSettings *helm.EnvSettings) error {
	c.PlanReleases(o.tags)

	if o.planValues {
		if err := c.PlanReleasesValues(o.dir); err != nil {
			return err
		}
	}

	if o.planRepos {
		c.PlanRepos()
	}

	if o.planReleases {
		if err := c.PlanManifests(o.dir, helmSettings); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) PlanRepos() {
	c.Repositories = repo.Plan(c.Releases, c.Repositories)
	names := make([]string, len(c.Repositories))
	for i, v := range c.Repositories {
		names[i] = v.Name
	}
	log.WithField("repositories", names).Info("🛠 Yml -> 🗄 repositories")
}

func (c *Config) PlanReleases(tags []string) {
	c.Releases = release.Plan(tags, c.Releases)
	names := make([]string, len(c.Releases))
	for i, v := range c.Releases {
		if feature.Dependencies {
			v.HandleDependencies(c.Releases)
		}
		names[i] = v.UniqName()
	}
	log.WithField("releases", names).Info("🛠 Yml -> 🛥 releases")
}

func (c *Config) PlanReleasesValues(dir string) error {
	return release.PlanValues(c.Releases, dir)
}

func (c *Config) PlanManifests(dir string, helmSettings *helm.EnvSettings) error {
	err := repo.Sync(c.Repositories, helmSettings)
	if err != nil {
		return err
	}

	for i, _ := range c.Releases {
		c.Releases[i].Options.DryRun = true
	}

	parallel := feature.Parallel
	feature.Parallel = false

	err = c.SyncReleases(dir + ".manifest/")

	feature.Parallel = parallel

	for i, _ := range c.Releases {
		c.Releases[i].Options.DryRun = false
	}

	return err
}
