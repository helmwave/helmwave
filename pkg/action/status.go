package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Status struct {
	plandir string
}

func (l *Status) Run() error {
	p := plan.New(l.plandir)
	return p.List()
}

func (l *Status) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List of deployed releases",
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
