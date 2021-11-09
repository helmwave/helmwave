package action

import (
	"os"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type DiffLive struct { //nolint:govet
	plandir string
	diff    *Diff
}

func (d *DiffLive) Run() error {
	p := plan.New(d.plandir)
	if err := p.Import(); err != nil {
		return err
	}
	if ok := p.IsManifestExist(); !ok {
		return os.ErrNotExist
	}

	p.DiffLive(d.diff.ShowSecret, d.diff.Wide)

	return nil
}

func (d *DiffLive) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "live",
		Usage: "plan 🆚 live",
		Flags: []cli.Flag{
			flagPlandir(&d.plandir),
		},
		Action: toCtx(d.Run),
	}
}
