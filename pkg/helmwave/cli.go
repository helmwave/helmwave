package helmwave

import (
	"fmt"
	"github.com/helmwave/helmwave/pkg/feature"
	"github.com/helmwave/helmwave/pkg/yml"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
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
		return cli.Exit(fmt.Sprintf("👻 Subcommand %q not found", ctx.Args().First()), 127)
	}

	opts := &yml.SavePlanOptions{}

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

	return c.plan(opts)
}

func (c *Config) plan(opts *yml.SavePlanOptions) error {
	if opts == nil {
		opts = &yml.SavePlanOptions{}
	}

	opts.File(c.PlanPath + PLANFILE).Dir(c.PlanPath)

	err := c.Tpl.Render()
	if err != nil {
		return err
	}

	err = yml.Read(c.Tpl.To, &c.Yml)
	if err != nil {
		return err
	}

	opts.Tags(c.Tags.Value())

	return c.Yml.SavePlan(opts, c.Helm)
}

func (c *Config) CliDeploy(ctx *cli.Context) error {
	opts := &yml.SavePlanOptions{}
	opts.PlanValues().PlanRepos().PlanValues()
	err := c.plan(opts)
	if err != nil {
		return err
	}

	if feature.Kubedog {
		log.Info("🐶 Kubedog enabled")
		return c.Yml.SyncWithKubedog(c.PlanPath+".manifest/", c.Helm, c.Kubedog)
	}

	return c.Yml.Sync(c.PlanPath+".manifest/", c.Helm)
}

func (c *Config) CliManifests(ctx *cli.Context) error {
	err := c.CliPlan(ctx)
	if err != nil {
		return err
	}

	return c.Yml.SyncFake(c.PlanPath+".manifest/", c.Helm)
}

func (c *Config) CliVersion(ctx *cli.Context) error {
	cli.ShowVersion(ctx)
	return nil
}

func (c *Config) CliStatus(ctx *cli.Context) error {
	err := c.plan(nil)
	if err != nil {
		return err
	}

	releases := ctx.Args().Slice()
	return c.Yml.Status(releases)
}

func (c *Config) CliList(ctx *cli.Context) error {
	opts := &yml.SavePlanOptions{}
	opts.PlanReleases()

	err := c.plan(opts)
	if err != nil {
		return err
	}

	return c.Yml.ListReleases()
}

func (c *Config) CliUninstall(ctx *cli.Context) error {
	opts := &yml.SavePlanOptions{}
	opts.PlanReleases()

	err := c.plan(opts)
	if err != nil {
		return err
	}

	return c.Yml.Uninstall()
}

func Command404(c *cli.Context, s string) {
	log.Errorf("👻 Command %q not found \n", s)
	os.Exit(127)
}
