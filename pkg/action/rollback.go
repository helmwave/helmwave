package action

import (
	"context"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

var _ Action = (*Rollback)(nil)

// Rollback is a struct for running 'rollback' command.
type Rollback struct {
	build     *Build
	dog       *kubedog.Config
	autoBuild bool
	revision  int
}

// Run is the main function for 'rollback' command.
func (i *Rollback) Run(ctx context.Context) error {
	if i.autoBuild {
		if err := i.build.Run(ctx); err != nil {
			return err
		}
	}
	p, err := plan.NewAndImport(ctx, i.build.planFS)
	if err != nil {
		return err
	}

	return p.Rollback(ctx, i.revision, i.dog)
}

// Cmd returns 'rollback' *cli.Command.
func (i *Rollback) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "rollback",
		Usage:  "⏮  rollback your plan",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
	}
}

// flags return flag set of CLI urfave.
func (i *Rollback) flags() []cli.Flag {
	// Init sub-structures
	i.build = &Build{}
	i.dog = &kubedog.Config{}

	self := []cli.Flag{
		flagAutoBuild(&i.autoBuild),
		&cli.IntFlag{
			Name:        "revision",
			Value:       -1,
			Usage:       "rollback all releases to this revision",
			Destination: &i.revision,
		},
	}

	self = append(self, flagsKubedog(i.dog)...)
	self = append(self, i.build.flags()...)

	return self
}
