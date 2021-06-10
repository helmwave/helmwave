package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Uninstall struct {
	Plandir string
}

func (i *Uninstall) Run() error {
	p := plan.New(i.Plandir)
	return p.Apply()
}

func (i *Uninstall) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "uninstall",
		Aliases: []string{"destroy", "delete", "del", "rm", "remove"},
		Usage:   "Delete deployed releases",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "plandir",
				Value:       ".helmwave/",
				Usage:       "Path to plandir",
				EnvVars:     []string{"HELMWAVE_PLANDIR"},
				Destination: &i.Plandir,
			},
		},
		Action: toCtx(i.Run),
	}
}
