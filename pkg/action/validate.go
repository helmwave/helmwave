package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

// Validate is struct for running 'validate' command.
type Validate struct {
	plandir string
}

// Run is main function for 'validate' command.
func (l *Validate) Run() error {
	p := plan.New(l.plandir)
	if err := p.Import(); err != nil {
		return err
	}

	return p.ValidateValues()
}

// Cmd returns 'validate' *cli.Command.
func (l *Validate) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "validate",
		Usage: "🛂 Validate your plan",
		Flags: []cli.Flag{
			flagPlandir(&l.plandir),
		},
		Action: toCtx(l.Run),
	}
}
