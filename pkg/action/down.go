package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Down struct {
	plandir string
}

func (i *Down) Run() error {
	p := plan.New(i.plandir)
	if err := p.Import(); err != nil {
		return err
	}

	return p.Destroy()
}

func (i *Down) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "down",
		Usage: "🔪 Delete all",
		Flags: []cli.Flag{
			flagPlandir(&i.plandir),
		},
		Action: toCtx(i.Run),
	}
}
