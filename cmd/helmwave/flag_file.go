package main

import "github.com/urfave/cli/v2"

var file = &cli.StringFlag{
	Name:        "file",
	Aliases:     []string{"f"},
	Value:       "helmwave.yml",
	Usage:       "Main yml file",
	EnvVars:     []string{"HELMWAVE_FILE", "HELMWAVE_YAML_FILE", "HELMWAVE_YML_FILE"},
	Destination: &app.Tpl.To,
}