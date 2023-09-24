package action

import (
	"context"
	"os"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

var _ Action = (*DiffLocal)(nil)

// DiffLocal is a struct for running 'diff plan' command.
type DiffLocal struct {
	diff    *Diff
	plan1FS plan.ImportFS
	plan2FS plan.ImportFS
}

// Run is the main function for 'diff plan' command.
func (d *DiffLocal) Run(ctx context.Context) error {
	// Plan 1
	plan1, err := plan.NewAndImport(ctx, d.plan1FS)
	if err != nil {
		return err
	}
	if ok := plan1.IsManifestExist(d.plan1FS); !ok {
		return os.ErrNotExist
	}

	// Plan 2
	plan2, err := plan.NewAndImport(ctx, d.plan2FS)
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
		&cli.GenericFlag{
			Name:        "plandir1",
			Destination: createGenericFS(&d.plan1FS, plan.Dir),
			DefaultText: getDefaultFSValue(plan.Dir),
			EnvVars:     []string{"HELMWAVE_PLANDIR_1", "HELMWAVE_PLANDIR"},
		},
		&cli.GenericFlag{
			Name:        "plandir2",
			Destination: createGenericFS(&d.plan2FS, plan.Dir),
			DefaultText: getDefaultFSValue(plan.Dir),
			EnvVars:     []string{"HELMWAVE_PLANDIR_2"},
		},
	}
}
