package action

import (
	"context"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

// Rollback is struct for running 'rollback' command.
type Rollback struct {
	build     *Build
	autoBuild bool
	revision  int
}

// Run is main function for 'rollback' command.
func (i *Rollback) Run(ctx context.Context) error {
	if i.autoBuild {
		if err := i.build.Run(ctx); err != nil {
			return err
		}
	}
	p, err := plan.NewAndImport(i.build.plandir)
	if err != nil {
		return err
	}

	return p.Rollback(i.revision)
}

// Cmd returns 'rollback' *cli.Command.
func (i *Rollback) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "rollback",
		Usage:  "⏮  Rollback your plan",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
	}
}

// flags return flag set of CLI urfave.
func (i *Rollback) flags() []cli.Flag {
	// Init sub-structures
	i.build = &Build{}

	self := []cli.Flag{
		flagAutoBuild(&i.autoBuild),
		&cli.IntFlag{
			Name:        "revision",
			Value:       -1,
			Usage:       "Rollback all releases to this revision",
			Destination: &i.revision,
		},
	}

	return append(self, i.build.flags()...)
}
