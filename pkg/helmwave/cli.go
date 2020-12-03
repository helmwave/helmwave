package helmwave

import (
	"github.com/urfave/cli/v2"
)

func (c *Config) Render(ctx *cli.Context) error {
	return c.RenderHelmWaveYml()
}

func (c *Config) Planfile(ctx *cli.Context) error {
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

func (c *Config) Repos(ctx *cli.Context) error {
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

func (c *Config) Deploy(ctx *cli.Context) error {
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

func (c *Config) UsePlan(ctx *cli.Context) error {
	c.InitPlanDirFile()

	err := c.ReadHelmWavePlan()
	if err != nil {
		return err
	}

	return c.SyncPlan()
}
