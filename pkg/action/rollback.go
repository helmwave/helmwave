package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

// Rollback is struct for running 'rollback' command.
type Rollback struct {
	build     *Build
	autoBuild bool
}

// Run is main function for 'rollback' command.
func (i *Rollback) Run() error {
	if i.autoBuild {
		if err := i.build.Run(); err != nil {
			return err
		}
	}
	p, err := plan.NewAndImport(i.build.plandir)
	if err != nil {
		return err
	}

	return p.Rollback()
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
		flagPlandir(&i.build.plandir),
	}

	return append(self, i.build.flags()...)
}
