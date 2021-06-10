package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Build struct {
	plandir string
}

func (i *Build) Run() error {
	p := plan.New(i.plandir)
	return p.Build()
}

func (i *Build) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Build plandir",
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
