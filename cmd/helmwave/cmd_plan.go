package main

import "github.com/urfave/cli/v2"

var plan = &cli.Command{
	Name:    "planfile",
	Aliases: []string{"plan"},
	Usage:   "📜 Generate planfile to plandir",
	Action:  app.CliPlan,
	Before:  app.InitApp,
	Flags: []cli.Flag{
		plandir,
		file,
		tags,
	},
}
