package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Status struct {
	plandir string
}

func (i *Status) Run() error {
	p := plan.New(i.plandir)
	if err := p.Import(); err != nil {
		return err
	}
	return p.List()
}

func (i *Status) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Status of deployed releases",
		Flags: []cli.Flag{
			flagPlandir(&i.plandir),
		},
		Action: toCtx(i.Run),
	}
}
