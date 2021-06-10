package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type List struct {
	Plandir string
}

func (l *List) Run() error {
	p := plan.New(l.Plandir)
	return p.List()
}

func (l *List) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List of deployed releases",
		Action:  toCtx(l.Run),
	}
}
