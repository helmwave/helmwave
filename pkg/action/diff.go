package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli/v2"
)

type Diff struct {
	plandir1 string
	plandir2 string
	action
}

func (d *Diff) Run() error {
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
		Usage: "📜 Diff 2 plans",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "plan B",
				Value:       ".helmwave/",
				Usage:       "Path to plandir A",
				EnvVars:     []string{"HELMWAVE_PLANDIR_A", "HELMWAVE_PLANDIR"},
				Destination: &d.plandir1,
			},
			&cli.StringFlag{
				Name:        "plan A",
				Usage:       "Path to plandir B",
				EnvVars:     []string{"HELMWAVE_PLANDIR_B"},
				Destination: &d.plandir1,
			},
		},
		Action: toCtx(d.Run),
	}
}
