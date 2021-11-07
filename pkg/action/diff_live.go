package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type DiffLive struct {
	plandir  string
	diffWide int
}

func (d *DiffLive) Run() error {
	p := plan.New(d.plandir)
	if err := p.Import(); err != nil {
		return err
	}

	p.DiffLive(d.diffWide)

	return nil
}

func (d *DiffLive) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "diff-live",
		Usage: "🆚 Differences between plan and live",
		Flags: []cli.Flag{
			flagPlandir(&d.plandir),
			flagDiffWide(&d.diffWide),
		},
		Action: toCtx(d.Run),
	}
}
