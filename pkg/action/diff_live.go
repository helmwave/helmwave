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
	planFS  plan.PlanImportFS
	plandir string
}

// Run is the main function for 'diff live' command.
func (d *DiffLive) Run(ctx context.Context) error {
	// TODO: get filesystems dynamically from args
	d.planFS = getBaseFS().(plan.PlanImportFS) //nolint:forcetypeassert

	p, err := plan.NewAndImport(ctx, d.planFS, d.plandir)
	if err != nil {
		return err
	}

	if ok := p.IsManifestExist(d.planFS); !ok {
		return os.ErrNotExist
	}

	p.DiffLive(ctx, d.planFS, d.diff.ShowSecret, d.diff.Wide, d.diff.ThreeWayMerge)

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
