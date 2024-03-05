package action

import (
	"context"
	"fmt"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/version"
	"github.com/urfave/cli/v2"
)

var _ Action = (*One)(nil)

const RELEASE_PLAN = ".helmwave_%s"

type One struct {
	release release.Config
	up      *Up
}

func (i *One) plandir() string {
	return fmt.Sprintf(RELEASE_PLAN, release.Default.Uniq())
}

func (i *One) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "one",
		Aliases: []string{"install", "upgrade"},
		Flags:   i.flags(),
		Action: func(c *cli.Context) error {
			i.release.SetName(c.Args().Get(0))
			i.release.SetChartName(c.Args().Get(1))
			return i.Run(getContextWithFlags(c))
		},
	}
}

func (i *One) flags() []cli.Flag {
	y := &Yml{}
	b := &Build{
		yml:     y,
		autoYml: true,
	}

	i.up = &Up{
		build:     b,
		autoBuild: true,
	}

	self := []cli.Flag{}

	self = append(self, i.up.flags()...)

	return nil
}

func (i *One) Run(ctx context.Context) (err error) {
	b := plan.NewBodyPillow()
	b.Project = i.plandir()
	b.Version = version.Version
	b.Releases = []release.Config{i.release}

	p := plan.New(i.plandir())

	return i.up.build.Run(ctx)

	//return i.up.Run(ctx)
}
