package main

import (
	"github.com/helmwave/helmwave/pkg/feature"
	"github.com/urfave/cli/v2"
	"time"
)


var flagsKubedog = []cli.Flag{
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
}