package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Install struct {
	plandir string
}

func (i *Install) Run() error {
	p := plan.New(i.plandir)
	return p.Apply()
}

func (i *Install) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "install",
		Aliases: []string{"apply", "sync", "deploy"},
		Usage:   "🛥 Deploy!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "plandir",
				Value:       ".helmwave/",
				Usage:       "Path to plandir",
				EnvVars:     []string{"HELMWAVE_PLANDIR"},
				Destination: &i.plandir,
			},
		},
		Action: toCtx(i.Run),
	}
}
