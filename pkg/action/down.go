package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

// Down is struct for running 'down' command.
type Down struct {
	build     *Build
	autoBuild bool
}

// Run is main function for 'down' command.
func (i *Down) Run() error {
	if i.autoBuild {
		if err := i.build.Run(); err != nil {
			return err
		}
	}

	p, err := plan.NewAndImport(i.build.plandir)
	if err != nil {
		return err
	}

	return p.Destroy()
}

// Cmd returns 'down' *cli.Command.
func (i *Down) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "down",
		Usage:  "🔪 Delete all",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
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
