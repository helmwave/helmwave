package main

import "github.com/urfave/cli/v2"

var plan = &cli.Command{
	Name:    "planfile",
	Aliases: []string{"plan"},
	Usage:   "📜 Generate planfile to plandir",
	Action:  app.CliPlan,
	Before:  app.InitApp,
	Subcommands: []*cli.Command{
		{
			Name:   "repos",
			Action: app.CliPlan,
		},
		{
			Name:   "releases",
			Action: app.CliPlan,
		},
		{
			Name:   "values",
			Action: app.CliPlan,
		},
	},
}