package action

import (
	"context"

	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var _ Action = (*Down)(nil)

// Down is a struct for running 'down' command.
type Down struct {
	build     *Build
	autoBuild bool
}

// Run is the main function for 'down' command.
func (i *Down) Run(ctx context.Context) error {
	if i.autoBuild {
		if err := i.build.Run(ctx); err != nil {
			return err
		}
	}

	if i.build.skipUnchanged {
		log.Warn("I can't use HELMWAVE_SKIP_UNCHANGED=true for down command. I've changed it to false")
		i.build.skipUnchanged = false
	}

	p, err := plan.NewAndImport(ctx, i.build.plandir)
	if err != nil {
		return err
	}

	return p.Down(ctx)
}

// Cmd returns 'down' *cli.Command.
func (i *Down) Cmd() *cli.Command {
	return &cli.Command{
		Name:     "down",
		Category: Step2,
		Usage:    "🔪 delete all",
		Flags:    i.flags(),
		Action:   toCtx(i.Run),
	}
}

// flags return flag set of CLI urfave.
func (i *Down) flags() []cli.Flag {
	// Init sub-structures
	i.build = &Build{}

	self := []cli.Flag{
		flagAutoBuild(&i.autoBuild),
	}

	return append(self, i.build.flags()...)
}
