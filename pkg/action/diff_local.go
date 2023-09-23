package action

import (
	"context"
	"os"

	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var _ Action = (*DiffLocal)(nil)

// DiffLocal is a struct for running 'diff plan' command.
type DiffLocal struct {
	diff     *Diff
	plan1FS  plan.PlanImportFS
	plan2FS  plan.PlanImportFS
	plandir1 string
	plandir2 string
}

// Run is the main function for 'diff plan' command.
func (d *DiffLocal) Run(ctx context.Context) error {
	// TODO: get filesystems dynamically from args
	d.plan1FS = getBaseFS().(plan.PlanImportFS) //nolint:forcetypeassert
	d.plan2FS = getBaseFS().(plan.PlanImportFS) //nolint:forcetypeassert

	if d.plandir1 == d.plandir2 {
		log.Warn(plan.ErrPlansAreTheSame)
	}

	// Plan 1
	plan1, err := plan.NewAndImport(ctx, d.plan1FS, d.plandir1)
	if err != nil {
		return err
	}
	if ok := plan1.IsManifestExist(d.plan1FS); !ok {
		return os.ErrNotExist
	}

	// Plan 2
	plan2, err := plan.NewAndImport(ctx, d.plan2FS, d.plandir2)
	if err != nil {
		return err
	}
	if ok := plan2.IsManifestExist(d.plan2FS); !ok {
		return os.ErrNotExist
	}

	plan1.DiffPlan(plan2, d.diff.ShowSecret, d.diff.Wide)

	return nil
}

// Cmd returns 'diff plan' *cli.Command.
func (d *DiffLocal) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "local",
		Aliases: []string{"plan"},
		Usage:   "plan1  🆚  plan2",
		Flags:   d.flags(),
		Action:  toCtx(d.Run),
	}
}

// flags return flag set of CLI urfave.
func (d *DiffLocal) flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "plandir1",
			Value:       ".helmwave/",
			Usage:       "path to plandir1",
			EnvVars:     []string{"HELMWAVE_PLANDIR_1", "HELMWAVE_PLANDIR"},
			Destination: &d.plandir1,
		},
		&cli.StringFlag{
			Name:        "plandir2",
			Value:       ".helmwave/",
			Usage:       "path to plandir2",
			EnvVars:     []string{"HELMWAVE_PLANDIR_2"},
			Destination: &d.plandir2,
		},
	}
}
