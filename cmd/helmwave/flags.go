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
			Name:        "parallel",
			Usage:       "Parallel mode",
			Value:       true,
			EnvVars:     []string{"HELMWAVE_PARALLEL"},
			Destination: &app.Parallel,
		},
		&cli.StringFlag{
			Name:        "log-format",
			Usage:       "You can set 'text' or 'json'",
			Value:       "text",
			EnvVars:     []string{"HELMWAVE_LOG_FORMAT"},
			Destination: &app.Logger.Format,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "You can set 'debug' or 'info' or 'warning' or 'panic' or 'fatal'",
			Value:       "info",
			EnvVars:     []string{"HELMWAVE_LOG_LEVEL", "HELMWAVE_LOG_LVL"},
			Destination: &app.Logger.Level,
		},
	}
}
