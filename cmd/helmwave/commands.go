package main

import (
	"github.com/urfave/cli/v2"
)

func commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "render",
			Usage:  "📄 Render tpl -> yml",
			Action: app.CliRender,
		},
		{
			Name:    "planfile",
			Aliases: []string{"plan"},
			Usage:   "📜 Generate planfile to plandir",
			Action:  app.CliPlanfile,
		},
		{
			Name:    "repos",
			Aliases: []string{"rep", "repo"},
			Usage:   "🗄 Sync repositories",
			Action:  app.CliRepos,
		},
		{
			Name:    "deploy",
			Aliases: []string{"apply", "sync", "release"},
			Usage:   "🛥 Deploy your helmwave!",
			Action:  app.CliDeploy,
		},
		{
			Name:      "help",
			Usage:     "🚑 Help me!",
			ArgsUsage: "[command]",
			Action:    help,
		},
		{
			Name:   "useplan",
			Usage:  "📜 -> 🛥 Deploy your helmwave from planfile!",
			Action: app.CliUsePlan,
		},
		{
			Name:    "manifests",
			Aliases: []string{"templates"},
			Usage:   "📜 -> 🛥 Show manifests from planfile!",
			Action:  app.CliManifest,
		},
	}

}

func help(c *cli.Context) error {
	args := c.Args()
	if args.Present() {
		return cli.ShowCommandHelp(c, args.First())
	}

	return cli.ShowAppHelp(c)
}
