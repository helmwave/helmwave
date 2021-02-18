package helmwave

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/zhilyaev/helmwave/pkg/yml"
)

func (c *Config) CliYml(ctx *cli.Context) error {
	return c.Tpl.Render()
}

func (c *Config) CliPlan(ctx *cli.Context) error {
	err := c.Tpl.Render()
	if err != err {
		return err
	}

	err = yml.Read(c.Tpl.To, &c.Yml)
	if err != err {
		return err
	}

	return c.Yml.SavePlan(c.PlanPath+PLANFILE, c.Tags.Value(), c.PlanPath)
}

func (c *Config) CliDeploy(ctx *cli.Context) error {
	err := c.CliPlan(ctx)
	if err != err {
		return err
	}

	if c.Kubedog {
		log.Info("Kubedog init")
		return c.Yml.SyncWithKubedog(c.PlanPath+".manifest/", c.Parallel, c.Helm)
	}

	return c.Yml.Sync(c.PlanPath+".manifest/", c.Parallel, c.Helm)
}

func (c *Config) CliManifests(ctx *cli.Context) error {
	err := c.CliPlan(ctx)
	if err != err {
		return err
	}

	return c.Yml.SyncFake(c.PlanPath+".manifest/", c.Parallel, c.Helm)
}
