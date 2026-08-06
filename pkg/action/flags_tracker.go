package action

import (
	"time"

	"github.com/helmwave/helmwave/pkg/tracker"
	"github.com/urfave/cli/v2"
)

func flagsTracker(dog *tracker.Config) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "kubedog",
			Usage:       "enable live tracking of deployed resources",
			Value:       false,
			Category:    CategoryKubedog,
			EnvVars:     EnvVars("KUBEDOG_ENABLED", "KUBEDOG"),
			Destination: &dog.Enabled,
		},
		&cli.DurationFlag{
			Name:        "kubedog-status-interval",
			Usage:       "grace period for the tracker to catch final resource states after deploy",
			Value:       5 * time.Second,
			Category:    CategoryKubedog,
			EnvVars:     EnvVars("KUBEDOG_STATUS_INTERVAL"),
			Destination: &dog.StatusInterval,
		},
		&cli.DurationFlag{
			Name:        "kubedog-start-delay",
			Usage:       "delay tracker start, don't make it too late",
			Value:       time.Second,
			Category:    CategoryKubedog,
			EnvVars:     EnvVars("KUBEDOG_START_DELAY"),
			Destination: &dog.StartDelay,
		},
		&cli.DurationFlag{
			Name:        "kubedog-timeout",
			Usage:       "timeout of resource tracking",
			Value:       5 * time.Minute,
			Category:    CategoryKubedog,
			EnvVars:     EnvVars("KUBEDOG_TIMEOUT"),
			Destination: &dog.Timeout,
		},
		&cli.BoolFlag{
			Name:        "kubedog-logs",
			Usage:       "stream logs of tracked pods",
			Value:       true,
			Category:    CategoryKubedog,
			EnvVars:     EnvVars("KUBEDOG_LOGS"),
			Destination: &dog.Logs,
		},
		&cli.IntFlag{
			Name:        "kubedog-log-width",
			Usage:       "set max width of streamed log lines, 0 to not trim",
			Value:       140,
			Category:    CategoryKubedog,
			EnvVars:     EnvVars("KUBEDOG_LOG_WIDTH"),
			Destination: &dog.LogWidth,
		},
		&cli.BoolFlag{
			Name:        "kubedog-track-all",
			Usage:       "track almost all resources, not only workloads",
			Value:       false,
			Category:    CategoryKubedog,
			EnvVars:     EnvVars("KUBEDOG_TRACK_ALL"),
			Destination: &dog.TrackGeneric,
		},
	}
}
