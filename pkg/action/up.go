package action

import (
	"context"

	"github.com/helmwave/helmwave/pkg/clictx"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var _ Action = (*Up)(nil)

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
	cliCtx := clictx.GetCLIFromContext(ctx)
	if cliCtx == nil {
		return
	}

	for _, buildFlag := range i.build.flags() {
		name := buildFlag.Names()[0]
		if cliCtx.IsSet(name) {
			log.WithField("flag", name).Warn("this flag is used by autobuild (--build) but autobuild is disabled")
		}
	}
}

// Cmd returns 'up' *cli.Command.
func (i *Up) Cmd() *cli.Command {
	return &cli.Command{
		Name:     "up",
		Category: Step2,
		Usage:    "🚢 apply your plan",
		Flags:    i.flags(),
		Action:   toCtx(i.Run),
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
			Name:        "progress",
			Usage:       "enable progress logs of helm (INFO log level)",
			Value:       false,
			Category:    "KUBEDOG",
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

	self = append(self, flagsKubedog(i.dog)...)
	self = append(self, i.build.flags()...)

	return self
}
