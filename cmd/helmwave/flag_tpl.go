package main

import "github.com/urfave/cli/v2"


var tpl  = &cli.StringFlag{
	Name:        "tpl",
	Value:       "helmwave.yml.tpl",
	Usage:       "Main tpl file",
	EnvVars:     []string{"HELMWAVE_TPL_FILE"},
	Destination: &app.Tpl.From,
}
