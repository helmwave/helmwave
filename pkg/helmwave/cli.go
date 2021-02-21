package helmwave

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/zhilyaev/helmwave/pkg/yml"
	"os"
	"fmt"
)

func (c *Config) InitApp(ctx *cli.Context) error {
	err := c.InitLogger()
	if err != nil {
		return err
	}

	c.InitPlan()
	return nil
}

func (c *Config) CliYml(ctx *cli.Context) error {
	return c.Tpl.Render()
}

func (c *Config) CliPlan(ctx *cli.Context) error {
	// We do not want any non-existing subcommands
	if ctx.Args().Present() {
		return cli.Exit(fmt.Sprintf("Subcommand %s not found", ctx.Args().First()), 127)
	}

	opts := &yml.SavePlanOptions{}
	opts.File(c.PlanPath+PLANFILE).Dir(c.PlanPath)

	switch ctx.Command.Name {
		case "repos":
			opts.PlanRepos()
		case "releases":
			opts.PlanReleases()
		case "values":
			opts.PlanValues()
		default:
			opts.PlanRepos().PlanReleases().PlanValues()
	}

	err := c.Tpl.Render()
	if err != nil {
		return err
	}

	err = yml.Read(c.Tpl.To, &c.Yml)
	if err != nil {
		return err
	}

	opts.Tags(c.Tags.Value())

	return c.Yml.SavePlan(opts)
}

func (c *Config) CliDeploy(ctx *cli.Context) error {
	err := c.CliPlan(ctx)
	if err != nil {
		return err
	}

	if c.Kubedog.Enabled {
		log.Info("🐶 Kubedog enabled")
		return c.Yml.SyncWithKubedog(c.PlanPath+".manifest/", c.Parallel, c.Helm, c.Kubedog.StatusInterval, c.Kubedog.Timeout)
	}

	return c.Yml.Sync(c.PlanPath+".manifest/", c.Parallel, c.Helm)
}

func (c *Config) CliManifests(ctx *cli.Context) error {
	err := c.CliPlan(ctx)
	if err != nil {
		return err
	}

	return c.Yml.SyncFake(c.PlanPath+".manifest/", c.Parallel, c.Helm)
}

func (c *Config) CliVersion(ctx *cli.Context) error {
	cli.ShowVersion(ctx)
	return nil
}

func Command404(c *cli.Context, s string) {
	log.Errorf("👻 Command %q not found \n", s)
	os.Exit(127)
}
