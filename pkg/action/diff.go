package action

import (
	"errors"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli/v2"
)

type Diff struct {
	plandir1 string
	plandir2 string
}

func (d *Diff) Run() error {
	if d.plandir1 == d.plandir2 {
		return errors.New("plan1 and plan2 are the same")
	}

	plan1 := plan.New(d.plandir1)
	if err := plan1.Import(); err != nil {
		return err
	}

	plan2 := plan.New(d.plandir2)
	if err := plan2.Import(); err != nil {
		return err
	}

	changelog, err := plan.Diff(plan1, plan2)
	log.Info(changelog)
	return err
}

func (d *Diff) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "diff",
		Usage: "Differences between plan1 and plan2",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "plandir1",
				Value:       ".helmwave/",
				Usage:       "Path to plandir1",
				EnvVars:     []string{"HELMWAVE_PLANDIR_1", "HELMWAVE_PLANDIR"},
				Destination: &d.plandir1,
			},
			&cli.StringFlag{
				Name:        "plandir2",
				Value:       ".helmwave/",
				Usage:       "Path to plandir2",
				EnvVars:     []string{"HELMWAVE_PLANDIR_2"},
				Destination: &d.plandir2,
			},
		},
		Action: toCtx(d.Run),
	}
}
