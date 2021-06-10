package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Install struct {
	Plandir string
}

func (i *Install) Run(c *cli.Context) error {
	p := plan.New(i.Plandir)
	return p.Apply()
}
