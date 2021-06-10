package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Build struct {
	plandir string
	tags    []string
}

func (i *Build) Run() error {
	p := plan.New(i.plandir)
	return p.Build()
}

func (i *Build) Cmd() *cli.Command {
	t := cli.StringSlice{}
	i.tags = t.Value()
	return &cli.Command{
		Name:  "build",
		Usage: "Build a plandir",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "plandir",
				Value:       ".helmwave/",
				Usage:       "Path to plandir",
				EnvVars:     []string{"HELMWAVE_PLANDIR"},
				Destination: &i.plandir,
			},
			&cli.StringSliceFlag{
				Name:        "tags",
				Aliases:     []string{"t"},
				Usage:       "It allows you choose releases for sync. Example: -t tag1 -t tag3,tag4",
				EnvVars:     []string{"HELMWAVE_TAGS"},
				Destination: &t,
			},
		},
		Action: toCtx(i.Run),
	}
}
