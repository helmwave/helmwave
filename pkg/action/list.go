package action

import (
	"context"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

// List is struct for running 'list' command.
type List struct {
	build     *Build
	autoBuild bool
}

// Run is main function for 'list' command.
func (l *List) Run(ctx context.Context) error {
	if l.autoBuild {
		if err := l.build.Run(ctx); err != nil {
			return err
		}
	}
	p, err := plan.NewAndImport(l.build.plandir)
	if err != nil {
		return err
	}

	return p.List()
}

// Cmd returns 'list' *cli.Command.
func (l *List) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "👀 List of deployed releases",
		Flags:   l.flags(),
		Action:  toCtx(l.Run),
	}
}

// flags return flag set of CLI urfave.
func (l *List) flags() []cli.Flag {
	// Init sub-structures
	l.build = &Build{}

	self := []cli.Flag{
		flagAutoBuild(&l.autoBuild),
	}

	return append(self, l.build.flags()...)
}
