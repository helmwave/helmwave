package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Rollback struct {
	plandir  string
	parallel bool
}

func (i *Rollback) Run() error {
	p := plan.New(i.plandir)
	return p.Rollback()
}

func (i *Rollback) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "Rollback",
		Usage: "Rollback your plandir",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "plandir",
				Value:       ".helmwave/",
				Usage:       "Path to plandir",
				EnvVars:     []string{"HELMWAVE_PLANDIR"},
				Destination: &i.plandir,
			},
			&cli.BoolFlag{
				Name:        "parallel",
				Usage:       "It allows you call `helm uninstall` in parallel mode ",
				Value:       true,
				EnvVars:     []string{"HELMWAVE_PARALLEL"},
				Destination: &i.parallel,
			},
		},
		Action: toCtx(i.Run),
	}
}
