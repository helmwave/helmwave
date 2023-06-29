package action

import (
	"context"

	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var _ Action = (*Graph)(nil)

// Graph is a struct for running 'graph' command.
type Graph struct {
	build     *Build
	autoBuild bool
}

// Run is the main function for 'status' command.
func (l *Graph) Run(ctx context.Context) error {
	if 1 == l.build.options.GraphWidth {
		log.Info("🔺it is not possible to turn off the graph in this command")

		return nil
	}

	old := l.build.options.GraphWidth
	if l.autoBuild {
		// Disable graph if it needs to build
		l.build.options.GraphWidth = 1
		if err := l.build.Run(ctx); err != nil {
			return err
		}
	}
	l.build.options.GraphWidth = old

	p, err := plan.NewAndImport(ctx, l.build.plandir)
	if err != nil {
		return err
	}

	log.Infof("show graph:\n%s", p.BuildGraphASCII(l.build.options.GraphWidth))

	return nil
}

// Cmd returns 'status' *cli.Command.
func (l *Graph) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "graph",
		Usage:  "show graph",
		Flags:  l.flags(),
		Action: toCtx(l.Run),
	}
}

// flags return flag set of CLI urfave.
func (l *Graph) flags() []cli.Flag {
	// Init sub-structures
	l.build = &Build{}

	self := []cli.Flag{
		flagAutoBuild(&l.autoBuild),
	}

	return append(self, l.build.flags()...)
}
