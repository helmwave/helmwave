package action

import (
	"context"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

// Status is a struct for running 'status' command.
type Status struct {
	build     *Build
	names     cli.StringSlice
	autoBuild bool
}

// Run is the main function for 'status' command.
func (l *Status) Run(ctx context.Context) error {
	if l.autoBuild {
		if err := l.build.Run(ctx); err != nil {
			return err
		}
	}
	p, err := plan.NewAndImport(ctx, l.build.plandir)
	if err != nil {
		return err
	}

	return p.Status(l.names.Value()...)
}

// Cmd returns 'status' *cli.Command.
func (l *Status) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "status",
		Usage:  "👁️status of deployed releases",
		Flags:  l.flags(),
		Action: toCtx(l.Run),
	}
}

// flags return flag set of CLI urfave.
func (l *Status) flags() []cli.Flag {
	// Init sub-structures
	l.build = &Build{}

	self := []cli.Flag{
		flagAutoBuild(&l.autoBuild),
	}

	return append(self, l.build.flags()...)
}
