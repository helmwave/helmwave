package helmwave

import (
	"github.com/urfave/cli/v2"
	"github.com/zhilyaev/helmwave/pkg/yml"
)

func (c *Config) CliYml(ctx *cli.Context) error {
	return c.Tpl.Render()
}

func (c *Config) CliPlan(ctx *cli.Context) error {
	err := yml.Read(c.Tpl.To, &c.Plan.Body)
	if err != err {
		return err
	}

	return c.Plan.Plan(c.Tags.Value(), c.PlanPath)
}

func (c *Config) CliDeploy(ctx *cli.Context) error {
	err := yml.Read(c.Tpl.To, &c.Plan.Body)
	if err != err {
		return err
	}

	return c.Plan.Plan(c.Tags.Value(), c.PlanPath)
}
