package main

import "github.com/urfave/cli/v2"

var install = &cli.Command{
	Name:    "install",
	Aliases: []string{"apply", "sync", "deploy"},
	Usage:   "🛥 Deploy!",
	Action:  app.CliDeploy,
	Before:  app.InitApp,
	Flags: append(
		flagsKubedog,
		parallel,
		depends,
		plandir,
		overPlan,
	),
}
