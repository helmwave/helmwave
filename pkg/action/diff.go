package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli/v2"
)

type Diff struct {
	Plandir1 string
	Plandir2 string
}

func (d *Diff) Run(c *cli.Context) error {
	plan1 := plan.New(d.Plandir1)
	if err := plan1.Import(); err != nil {
		return err
	}

	plan2 := plan.New(d.Plandir2)
	if err := plan2.Import(); err != nil {
		return err
	}

	changelog, err := plan.Diff(plan1, plan2)
	log.Info(changelog)
	return err
}
