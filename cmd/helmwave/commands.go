package main

import (
	"github.com/urfave/cli/v2"
)

func commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "render",
			Aliases: []string{"r"},
			Usage:   "📄 Render tpl -> yml",
			Action:  app.Render,
		},
		{
			Name:    "planfile",
			Aliases: []string{"p", "plan"},
			Usage:   "📜 Generate planfile",
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
			Aliases: []string{"d", "apply", "sync", "release"},
			Usage:   "🛥 Deploy your helmwave!",
			Action:  app.SyncReleases,
		},
		{
			Name:      "help",
			Usage:     "🚑 Help me!",
			Aliases:   []string{"h"},
			ArgsUsage: "[command]",
			Action:    help,
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
