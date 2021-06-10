package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Uninstall struct {
	plandir  string
	parallel bool
}

func (i *Uninstall) Run() error {
	p := plan.New(i.plandir)
	return p.Apply()
}

func (i *Uninstall) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "uninstall",
		Aliases: []string{"destroy", "delete", "del", "rm", "remove"},
		Usage:   "Delete all",
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
