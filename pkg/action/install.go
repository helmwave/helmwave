package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Install struct {
	plandir  string
	parallel bool
	kubedog  bool
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
			&cli.BoolFlag{
				Name:        "parallel",
				Usage:       "It allows you call `helm install` in parallel mode ",
				Value:       true,
				EnvVars:     []string{"HELMWAVE_PARALLEL"},
				Destination: &i.parallel,
			},
			&cli.BoolFlag{
				Name:        "kubedog",
				Usage:       "Enable/Disable kubedog",
				Value:       false,
				EnvVars:     []string{"HELMWAVE_KUBEDOG"},
				Destination: &i.kubedog,
			},
		},
		Action: toCtx(i.Run),
	}
}
