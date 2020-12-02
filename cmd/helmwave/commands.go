package main

import (
	"github.com/urfave/cli/v2"
)

func commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "render",
			Usage:  "📄 Render tpl -> yml",
			Action: app.Render,
		},
		{
			Name:    "planfile",
			Aliases: []string{"plan"},
			Usage:   "📜 Generate planfile to plandir",
			Action:  app.Planfile,
		},
		{
			Name:    "repos",
			Aliases: []string{"rep", "repo"},
			Usage:   "🗄 Sync repositories",
			Action:  app.SyncRepos,
		},
		{
			Name:    "deploy",
			Aliases: []string{"apply", "sync", "release"},
			Usage:   "🛥 Deploy your helmwave!",
			Action:  app.SyncReleases,
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
			Action: app.UsePlan,
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
