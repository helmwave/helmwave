package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

// Status is struct for running 'status' command.
type Status struct {
	plandir string
	names   cli.StringSlice
}

// Run is main function for 'status' command.
func (l *Status) Run() error {
	p := plan.New(l.plandir)
	if err := p.Import(); err != nil {
		return err
	}

	return p.Status(l.names.Value())
}

// Cmd returns 'status' *cli.Command.
func (l *Status) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "👁️ Status of deployed releases",
		Flags: []cli.Flag{
			flagPlandir(&l.plandir),
		},
		Action: toCtx(l.Run),
	}
}
