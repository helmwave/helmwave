package action

import (
	"context"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

var _ Action = (*Validate)(nil)

// Validate is a struct for running 'validate' command.
type Validate struct {
	planFS  plan.PlanImportFS
	plandir string
}

// Run is the main function for 'validate' command.
func (l *Validate) Run(ctx context.Context) error {
	// TODO: get filesystems dynamically from args
	l.planFS = getBaseFS().(plan.PlanImportFS) //nolint:forcetypeassert

	p, err := plan.NewAndImport(ctx, l.planFS, l.plandir)
	if err != nil {
		return err
	}

	return p.ValidateValuesImport(l.planFS)
}

// Cmd returns 'validate' *cli.Command.
func (l *Validate) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "validate",
		Usage:  "🛂 validate your plan",
		Flags:  l.flags(),
		Action: toCtx(l.Run),
	}
}

// flags return flag set of CLI urfave.
func (l *Validate) flags() []cli.Flag {
	return []cli.Flag{
		flagPlandir(&l.plandir),
	}
}
