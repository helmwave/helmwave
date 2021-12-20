package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type List struct {
	plandir string
}

func (l *List) Run() error {
	p := plan.New(l.plandir)
	if err := p.Import(); err != nil {
		return err
	}

	return p.List()
}

func (l *List) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "👀 List of deployed releases",
		Flags: []cli.Flag{
			flagPlandir(&l.plandir),
		},
		Action: toCtx(l.Run),
	}
}
