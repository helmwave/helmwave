package main

import (
	"github.com/urfave/cli/v2"
	"github.com/zhilyaev/helmwave/pkg/helmwave"
)

func flags(app *helmwave.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "tpl",
			Value:       "helmwave.yml.tpl",
			Usage:       "Main tpl file",
			EnvVars:     []string{"HELMWAVE_TPL_FILE"},
			Destination: &app.Tpl.File,
		},
		&cli.StringFlag{
			Name:        "file",
			Aliases:     []string{"f"},
			Value:       "helmwave.yml",
			Usage:       "Main yml file",
			EnvVars:     []string{"HELMWAVE_FILE", "HELMWAVE_YAML_FILE", "HELMWAVE_YML_FILE"},
			Destination: &app.Yml.File,
		},
		&cli.StringFlag{
			Name:        "planfile",
			Aliases:     []string{"p"},
			Value:       "helmwave.plan",
			EnvVars:     []string{"HELMWAVE_PLANFILE"},
			Destination: &app.Plan.File,
		},
		&cli.StringSliceFlag{
			Name:        "tags",
			Aliases:     []string{"t"},
			Usage:       "Chose tags: -t tag1 -t tag3,tag4",
			EnvVars:     []string{"HELMWAVE_TAGS"},
			Destination: &app.Tags,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Usage:       "Debug helmwave",
			Aliases:     []string{"d"},
			Value:       false,
			EnvVars:     []string{"HELMWAVE_DEBUG"},
			Destination: &app.Debug,
		},
		&cli.BoolFlag{
			Name:        "parallel",
			Usage:       "Parallel mode",
			Value:       true,
			EnvVars:     []string{"HELMWAVE_PARALLEL"},
			Destination: &app.Parallel,
		},
	}
}
