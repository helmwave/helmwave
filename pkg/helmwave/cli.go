package helmwave

import (
	"github.com/urfave/cli/v2"
)

func (c *Config) CliRender(ctx *cli.Context) error {
	return c.RenderHelmWaveYml()
}

func (c *Config) CliPlanfile(ctx *cli.Context) error {
	err := c.RenderHelmWaveYml()
	if err != nil {
		return err
	}

	err = c.ReadHelmWaveYml()
	if err != nil {
		return err
	}

	return c.GenPlanfile()
}

func (c *Config) CliRepos(ctx *cli.Context) error {
	err := c.RenderHelmWaveYml()
	if err != nil {
		return err
	}

	err = c.ReadHelmWaveYml()
	if err != nil {
		return err
	}

	err = c.GenPlanfile()
	if err != nil {
		return err
	}

	return c.SyncPlanRepos()
}

func (c *Config) CliDeploy(ctx *cli.Context) error {
	err := c.RenderHelmWaveYml()
	if err != nil {
		return err
	}

	err = c.ReadHelmWaveYml()
	if err != nil {
		return err
	}

	err = c.GenPlanfile()
	if err != nil {
		return err
	}

	return c.SyncPlan()
}

func (c *Config) CliUsePlan(ctx *cli.Context) error {
	c.InitPlanDirFile()

	err := c.ReadHelmWavePlan()
	if err != nil {
		return err
	}

	return c.SyncPlan()
}

func (c *Config) CliManifest(ctx *cli.Context) error {
	c.InitPlanDirFile()
	err := c.ReadHelmWavePlan()
	if err != nil {
		return err
	}

	for i, _ := range c.Plan.Body.Releases {
		c.Plan.Body.Releases[i].Options.DryRun = true

	}

	return c.SyncPlan()
}
