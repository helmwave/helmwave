package main

import (
	"github.com/urfave/cli/v2"
)


var flagsLog = []cli.Flag{
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
	&cli.IntFlag{
		Name:        "kubedog-log-width",
		Usage:       "Set kubedog max log line width",
		Value:       140,
		EnvVars:     []string{"HELMWAVE_KUBEDOG_LOG_WIDTH"},
		Destination: &app.Logger.Width,
	},
}