package action

import (
	"context"
	"time"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Up is a struct for running 'up' command.
type Up struct {
	build     *Build
	dog       *kubedog.Config
	autoBuild bool
}

// Run is the main function for 'up' command.
func (i *Up) Run(ctx context.Context) error {
	if i.autoBuild {
		if err := i.build.Run(ctx); err != nil {
			return err
		}
	} else {
		i.warnOnBuildFlags(ctx)
	}

	p, err := plan.NewAndImport(ctx, i.build.plandir)
	if err != nil {
		return err
	}

	p.Logger().Info("🏗 Plan")

	return p.Up(ctx, i.dog)
}

func (i *Up) warnOnBuildFlags(ctx context.Context) {
	cliCtx, ok := ctx.Value("cli").(*cli.Context)
	if ok && cliCtx != nil {
		for _, buildFlag := range i.build.flags() {
			name := buildFlag.Names()[0]
			if cliCtx.IsSet(name) {
				log.WithField("flag", name).Warn("this flag is used by autobuild (--build) but autobuild is disabled")
			}
		}
	}
}

// Cmd returns 'up' *cli.Command.
func (i *Up) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "up",
		Usage:  "🚢 apply your plan",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
	}
}

// flags return flag set of CLI urfave.
func (i *Up) flags() []cli.Flag {
	// Init sub-structures
	i.dog = &kubedog.Config{}
	i.build = &Build{}

	self := []cli.Flag{
		flagAutoBuild(&i.autoBuild),
		&cli.BoolFlag{
			Name:        "kubedog",
			Usage:       "enable/disable kubedog",
			Value:       false,
			EnvVars:     []string{"HELMWAVE_KUBEDOG_ENABLED", "HELMWAVE_KUBEDOG"},
			Destination: &i.dog.Enabled,
		},
		&cli.DurationFlag{
			Name:        "kubedog-status-interval",
			Usage:       "interval of kubedog status messages",
			Value:       5 * time.Second,
			EnvVars:     []string{"HELMWAVE_KUBEDOG_STATUS_INTERVAL"},
			Destination: &i.dog.StatusInterval,
		},
		&cli.DurationFlag{
			Name:        "kubedog-start-delay",
			Usage:       "delay kubedog start, don't make it too late",
			Value:       time.Second,
			EnvVars:     []string{"HELMWAVE_KUBEDOG_START_DELAY"},
			Destination: &i.dog.StartDelay,
		},
		&cli.DurationFlag{
			Name:        "kubedog-timeout",
			Usage:       "timeout of kubedog multitrackers",
			Value:       5 * time.Minute,
			EnvVars:     []string{"HELMWAVE_KUBEDOG_TIMEOUT"},
			Destination: &i.dog.Timeout,
		},
		&cli.IntFlag{
			Name:        "kubedog-log-width",
			Usage:       "set kubedog max log line width",
			Value:       140,
			EnvVars:     []string{"HELMWAVE_KUBEDOG_LOG_WIDTH"},
			Destination: &i.dog.LogWidth,
		},
		&cli.BoolFlag{
			Name:        "progress",
			Usage:       "enable progress logs of helm (INFO log level)",
			Value:       false,
			EnvVars:     []string{"HELMWAVE_PROGRESS"},
			Destination: &helper.Helm.Debug,
		},
		&cli.IntFlag{
			Name:    "parallel-limit",
			Usage:   "limit amount of parallel releases",
			EnvVars: []string{"HELMWAVE_PARALLEL_LIMIT"},
			Value:   0,
		},
	}

	return append(self, i.build.flags()...)
}
