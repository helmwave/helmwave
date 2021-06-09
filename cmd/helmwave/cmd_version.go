package main

import "github.com/urfave/cli/v2"

var version = &cli.Command{
	Name:   "version",
	Usage:  "Print helmwave version",
	Action: app.CliVersion,
}
