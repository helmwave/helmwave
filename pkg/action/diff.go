package action

import (
	"errors"
	"os"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Diff struct {
	plandir1, plandir2 string
	diffWide           int
}

func (d *Diff) Run() error {
	if d.plandir1 == d.plandir2 {
		return errors.New("plan1 and plan2 are the same")
	}

	plan1 := plan.New(d.plandir1)
	if err := plan1.Import(); err != nil {
		return err
	}
	if ok := plan1.IsManifestExist(); !ok {
		return os.ErrNotExist
	}

	plan2 := plan.New(d.plandir2)
	if err := plan2.Import(); err != nil {
		return err
	}
	if ok := plan2.IsManifestExist(); !ok {
		return os.ErrNotExist
	}

	plan1.Diff(plan2, d.diffWide)

	return nil
}

func (d *Diff) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "diff",
		Usage:   "🆚 Differences between plan1 and plan2",
		Aliases: []string{"vs"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "plandir1",
				Value:       ".helmwave/",
				Usage:       "Path file plandir1",
				EnvVars:     []string{"HELMWAVE_PLANDIR_1", "HELMWAVE_PLANDIR"},
				Destination: &d.plandir1,
			},
			&cli.StringFlag{
				Name:        "plandir2",
				Value:       ".helmwave/",
				Usage:       "Path file plandir2",
				EnvVars:     []string{"HELMWAVE_PLANDIR_2"},
				Destination: &d.plandir2,
			},
			flagDiffWide(&d.diffWide),
		},
		Action: toCtx(d.Run),
	}
}
