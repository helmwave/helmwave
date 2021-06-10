package plan

import (
	helm "helm.sh/helm/v3/pkg/cli"
)

func (o *SaveOptions) Plan(c *yml.Config, helmSettings *helm.EnvSettings) error {
	c.PlanReleases(o.tags)

	if o.withValues {
		if err := c.PlanReleasesValues(o.dir); err != nil {
			return err
		}
	}

	if o.withRepos {
		c.PlanRepos()
	}

	if o.withReleases {
		if err := c.PlanManifests(o.dir, helmSettings); err != nil {
			return err
		}
	}

	return nil
}

func (o *SaveOptions) File(file string) *SaveOptions {
	o.file = file
	return o
}

func (o *SaveOptions) Tags(tags []string) *SaveOptions {
	o.tags = tags
	return o
}

func (o *SaveOptions) Dir(dir string) *SaveOptions {
	o.dir = dir
	return o
}

func (o *SaveOptions) PlanReleases() *SaveOptions {
	o.withReleases = true
	return o.PlanValues().PlanRepos()
}

func (o *SaveOptions) PlanRepos() *SaveOptions {
	o.withRepos = true
	return o
}

func (o *SaveOptions) PlanValues() *SaveOptions {
	o.withValues = true
	return o
}

func (o *SaveOptions) GetFile() string {
	return o.file
}
