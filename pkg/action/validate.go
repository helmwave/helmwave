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
		Name:  "validate",
		Usage: "🛂 Validate your plan",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "plandir",
				Value:       ".helmwave/",
				Usage:       "Path to plandir",
				EnvVars:     []string{"HELMWAVE_PLANDIR"},
				Destination: &l.plandir,
			},
		},
		Action: toCtx(l.Run),
	}
}
