package main

import "github.com/urfave/cli/v2"

var tags = &cli.StringSliceFlag{
	Name:        "tags",
	Aliases:     []string{"t"},
	Usage:       "It allows you choose releases. Example: -t tag1 -t tag3,tag4",
	EnvVars:     []string{"HELMWAVE_TAGS"},
	Destination: &app.Tags,
}