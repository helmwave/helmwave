package main

import "github.com/urfave/cli/v2"

var manifest = &cli.Command{
	Name:    "manifest",
	Usage:   "🛥 Fake Deploy",
	Action:  app.CliManifests,
	Before:  app.InitApp,
	Flags: []cli.Flag{
		parallel,
	},
}