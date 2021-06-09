package main

import "github.com/urfave/cli/v2"

var yml = &cli.Command	{
	Name:   "yml",
	Usage:  "📄 Render helmwave.yml.tpl -> helmwave.yml",
	Action: app.CliYml,
	Before: app.InitApp,
	Flags: []cli.Flag{
		tpl,
		file,
	},
}