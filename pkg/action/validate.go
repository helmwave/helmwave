package action

import (
	"context"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

var _ Action = (*Validate)(nil)

// Validate is a struct for running 'validate' command.
type Validate struct {
	plandir string
}

// Run is the main function for 'validate' command.
func (l *Validate) Run(ctx context.Context) error {
	p, err := plan.NewAndImport(ctx, l.plandir)
	if err != nil {
		return err
	}

	return p.ValidateValuesImport()
}

// Cmd returns 'validate' *cli.Command.
func (l *Validate) Cmd() *cli.Command {
	return &cli.Command{
		Name:     "validate",
		Category: Step1,
		Usage:    "🛂 validate your plan",
		Flags:    l.flags(),
		Action:   toCtx(l.Run),
	}
}

// flags return flag set of CLI urfave.
func (l *Validate) flags() []cli.Flag {
	return []cli.Flag{
		flagPlandir(&l.plandir),
	}
}
