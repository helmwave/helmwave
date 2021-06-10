package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Manifest struct {
	plandir string
}

func (i *Manifest) Run() error {
	p := plan.New(i.plandir)
	return p.Apply()
}

func (i *Manifest) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "Manifest",
		Usage: "Makes manifests",
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
