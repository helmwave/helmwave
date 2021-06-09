package main

import "github.com/urfave/cli/v2"

var list = &cli.Command		{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List of deployed releases",
	Action:  app.CliList,
	Before:  app.InitApp,
	Flags: []cli.Flag{
		plandir,
	},
}
