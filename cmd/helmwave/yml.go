package main

import (
	"github.com/helmwave/helmwave/pkg/action"
	"github.com/urfave/cli/v2"
)

var aYml = &action.Yml{}

var yml = &cli.Command	{
	Name:   "yml",
	Usage:  "📄 Render helmwave.yml.tpl -> helmwave.yml",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "tpl",
			Value:       "helmwave.yml.tpl",
			Usage:       "Main tpl file",
			EnvVars:     []string{"HELMWAVE_TPL_FILE"},
			Destination: &aYml.From,
		},
		&cli.StringFlag{
			Name:        "file",
			Aliases:     []string{"f"},
			Value:       "helmwave.yml",
			Usage:       "Main yml file",
			EnvVars:     []string{"HELMWAVE_FILE", "HELMWAVE_YAML_FILE", "HELMWAVE_YML_FILE"},
			Destination: &aYml.To,
		},
	},
	Action: func(c *cli.Context) error {
		return aYml.Run()
	},
}

