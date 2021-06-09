package main

import "github.com/urfave/cli/v2"

var uninstall = &cli.Command	{
	Name:    "uninstall",
	Aliases: []string{"destroy", "delete", "del", "rm", "remove"},
	Usage:   "Delete deployed releases",
	Action:  app.CliUninstall,
	Before:  app.InitApp,
	Flags: []cli.Flag{
		parallel,
		plandir,
	},
}
