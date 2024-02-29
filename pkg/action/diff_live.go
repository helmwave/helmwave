package action

import (
	"context"
	"os"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

var _ Action = (*DiffLive)(nil)

// DiffLive is a struct for running 'diff live' command.
type DiffLive struct {
	diff    *Diff
	plandir string
}

// Run is the main function for 'diff live' command.
func (d *DiffLive) Run(ctx context.Context) error {
	p, err := plan.NewAndImport(ctx, d.plandir)
	if err != nil {
		return err
	}

	if ok := p.IsManifestExist(); !ok {
		return os.ErrNotExist
	}

	p.DiffLive(ctx, d.diff.Options, d.diff.ThreeWayMerge)

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
		flagDiffThreeWayMerge(&d.diff.ThreeWayMerge),
	}
}
