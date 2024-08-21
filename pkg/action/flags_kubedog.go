package action

import (
	"time"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/urfave/cli/v2"
)

func flagsKubedog(dog *kubedog.Config) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "kubedog",
			Usage:       "enable/disable kubedog",
			Value:       false,
			Category:    "KUBEDOG",
			EnvVars:     EnvVars("KUBEDOG_ENABLED", "KUBEDOG"),
			Destination: &dog.Enabled,
		},
		&cli.DurationFlag{
			Name:        "kubedog-status-interval",
			Usage:       "interval of kubedog status messages: set -1s to stop showing status progress",
			Value:       5 * time.Second,
			Category:    "KUBEDOG",
			EnvVars:     EnvVars("KUBEDOG_STATUS_INTERVAL"),
			Destination: &dog.StatusInterval,
		},
		&cli.DurationFlag{
			Name:        "kubedog-start-delay",
			Usage:       "delay kubedog start, don't make it too late",
			Value:       time.Second,
			Category:    "KUBEDOG",
			EnvVars:     EnvVars("KUBEDOG_START_DELAY"),
			Destination: &dog.StartDelay,
		},
		&cli.DurationFlag{
			Name:        "kubedog-timeout",
			Usage:       "timeout of kubedog multitrackers",
			Value:       5 * time.Minute,
			Category:    "KUBEDOG",
			EnvVars:     EnvVars("KUBEDOG_TIMEOUT"),
			Destination: &dog.Timeout,
		},
		&cli.IntFlag{
			Name:        "kubedog-log-width",
			Usage:       "set kubedog max log line width",
			Value:       140,
			Category:    "KUBEDOG",
			EnvVars:     EnvVars("KUBEDOG_LOG_WIDTH"),
			Destination: &dog.LogWidth,
		},
		&cli.BoolFlag{
			Name:        "kubedog-track-all",
			Usage:       "track almost all resources, experimental",
			Value:       false,
			Category:    "KUBEDOG",
			EnvVars:     EnvVars("KUBEDOG_TRACK_ALL"),
			Destination: &dog.TrackGeneric,
		},
	}
}
