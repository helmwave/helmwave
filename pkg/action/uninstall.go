package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Uninstall struct {
	Plandir string
}

func (i *Uninstall) Run(c *cli.Context) error {
	p := plan.New(i.Plandir)
	return p.Apply()
}
