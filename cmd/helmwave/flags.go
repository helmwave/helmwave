package main

import (
	"github.com/helmwave/helmwave/pkg/feature"
	"github.com/helmwave/helmwave/pkg/helmwave"
	"github.com/urfave/cli/v2"
	"time"
)

func flags(app *helmwave.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "tpl",
			Value:       "helmwave.yml.tpl",
			Usage:       "Main tpl file",
			EnvVars:     []string{"HELMWAVE_TPL_FILE"},
			Destination: &app.Tpl.From,
		},
		&cli.StringFlag{
			Name:        "file",
			Aliases:     []string{"f"},
			Value:       "helmwave.yml",
			Usage:       "Main yml file",
			EnvVars:     []string{"HELMWAVE_FILE", "HELMWAVE_YAML_FILE", "HELMWAVE_YML_FILE"},
			Destination: &app.Tpl.To,
		},
		&cli.StringFlag{
			Name:        "plan-dir",
			Value:       ".helmwave/",
			Usage:       "It keeps your state via planfile",
			EnvVars:     []string{"HELMWAVE_PLAN_DIR"},
			Destination: &app.PlanPath,
		},
		&cli.StringSliceFlag{
			Name:        "tags",
			Aliases:     []string{"t"},
			Usage:       "It allows you choose releases for sync. Example: -t tag1 -t tag3,tag4",
			EnvVars:     []string{"HELMWAVE_TAGS"},
			Destination: &app.Tags,
		},
		&cli.BoolFlag{
			Name:        "parallel",
			Usage:       "It allows you call `helm install` in parallel mode ",
			Value:       true,
			EnvVars:     []string{"HELMWAVE_PARALLEL"},
			Destination: &feature.Parallel,
		},
		//&cli.BoolFlag{
		//	Name:        "force",
		//	Usage:       "It allows you call `helm install` in parallel mode ",
		//	Value:       true,
		//	EnvVars:     []string{"HELMWAVE_FORCE"},
		//	Destination: &app.Force,
		//},
		//
		//		LOGGER
		//
		&cli.StringFlag{
			Name:        "log-format",
			Usage:       "You can set: [ text | json | pad | emoji ]",
			Value:       "emoji",
			EnvVars:     []string{"HELMWAVE_LOG_FORMAT"},
			Destination: &app.Logger.Format,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "You can set: [ debug | info | warn  | fatal | panic | trace ]",
			Value:       "info",
			EnvVars:     []string{"HELMWAVE_LOG_LEVEL", "HELMWAVE_LOG_LVL"},
			Destination: &app.Logger.Level,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Force color",
			Value:       true,
			EnvVars:     []string{"HELMWAVE_LOG_COLOR"},
			Destination: &app.Logger.Color,
		},
		//
		// Kubedog Config
		&cli.BoolFlag{
			Name:        "kubedog",
			Usage:       "Enable/Disable kubedog",
			Value:       false,
			EnvVars:     []string{"HELMWAVE_KUBEDOG", "HELMWAVE_KUBEDOG_ENABLED"},
			Destination: &feature.Kubedog,
		},
		&cli.DurationFlag{
			Name:        "kubedog-status-interval",
			Usage:       "Interval of kubedog status messages",
			Value:       5 * time.Second,
			EnvVars:     []string{"HELMWAVE_KUBEDOG_STATUS_INTERVAL"},
			Destination: &app.Kubedog.StatusInterval,
		},
		&cli.DurationFlag{
			Name:        "kubedog-start-delay",
			Usage:       "Delay kubedog start, don't make it too late",
			Value:       time.Second,
			EnvVars:     []string{"HELMWAVE_KUBEDOG_START_DELAY"},
			Destination: &app.Kubedog.StartDelay,
		},
		&cli.DurationFlag{
			Name:        "kubedog-timeout",
			Usage:       "Timout of kubedog multitrackers",
			Value:       5 * time.Minute,
			EnvVars:     []string{"HELMWAVE_KUBEDOG_TIMEOUT"},
			Destination: &app.Kubedog.Timeout,
		},
		&cli.IntFlag{
			Name:        "kubedog-log-width",
			Usage:       "Set kubedog max log line width",
			Value:       140,
			EnvVars:     []string{"HELMWAVE_KUBEDOG_LOG_WIDTH"},
			Destination: &app.Logger.Width,
		},
		&cli.BoolFlag{
			Name:        "enable-dependencies",
			Usage:       "Enable dependencies",
			Value:       false,
			EnvVars:     []string{"HELMWAVE_ENABLE_DEPENDENCIES"},
			Destination: &feature.Dependencies,
		},
		&cli.BoolFlag{
			Name:        "plan-dependencies",
			Usage:       "Automatically add dependencies to plan",
			Value:       false,
			EnvVars:     []string{"HELMWAVE_PLAN_DEPENDENCIES"},
			Destination: &feature.PlanDependencies,
		},
	}
}
