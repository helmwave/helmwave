package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Uninstall struct {
	plandir string
	//parallel bool
}

func (i *Uninstall) Run() error {
	p := plan.New(i.plandir)
	if err := p.Import(); err != nil {
		return err
	}
	return p.Apply()
}

func (i *Uninstall) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "uninstall",
		Aliases: []string{"destroy", "delete", "del", "rm", "remove"},
		Usage:   "🔪 Delete all",
		Flags: []cli.Flag{
			flagPlandir(&i.plandir),
			//flagParallel(&i.parallel),
		},
		Action: toCtx(i.Run),
	}
}
