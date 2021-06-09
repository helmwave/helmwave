package main

import "github.com/urfave/cli/v2"

var status = &cli.Command{
	Name:    "status",
	Usage:   "Show status",
	Action:  app.CliDeploy,
	Before:  app.InitApp,
	Flags: []cli.Flag{
		plandir,
	},
}
