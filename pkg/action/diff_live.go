package action

import (
	"os"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

// DiffLive is struct for running 'diff live' command.
type DiffLive struct {
	diff    *Diff
	plandir string
}

// Run is main function for 'diff live' command.
func (d *DiffLive) Run() error {
	p, err := plan.NewAndImport(d.plandir)
	if err != nil {
		return err
	}

	if ok := p.IsManifestExist(); !ok {
		return os.ErrNotExist
	}

	p.DiffLive(d.diff.ShowSecret, d.diff.Wide)

	return nil
}

// Cmd returns 'diff live' *cli.Command.
func (d *DiffLive) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "live",
		Usage:  "plan 🆚 live",
		Flags:  d.flags(),
		Action: toCtx(d.Run),
	}
}

// flags return flag set of CLI urfave.
func (d *DiffLive) flags() []cli.Flag {
	return []cli.Flag{
		flagPlandir(&d.plandir),
	}
}
