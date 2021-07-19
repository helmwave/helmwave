package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Down struct {
	plandir string
	// parallel bool
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
		Name:    "down",
		Aliases: []string{"uninstall", "destroy", "delete", "del", "rm", "remove"},
		Usage:   "🔪 Delete all",
		Flags: []cli.Flag{
			flagPlandir(&i.plandir),
			// flagParallel(&i.parallel),
		},
		Action: toCtx(i.Run),
	}
}
