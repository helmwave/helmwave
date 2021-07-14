package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Validate struct {
	plandir string
}

func (l *Validate) Run() error {
	p := plan.New(l.plandir)
	if err := p.Import(); err != nil {
		return err
	}
	return p.ValidateValues()
}

func (l *Validate) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "validate",
		Aliases: []string{"check", "lint"},
		Usage:   "🛂 Validate your plan",
		Flags: []cli.Flag{
			flagPlandir(&l.plandir),
		},
		Action: toCtx(l.Run),
	}
}
