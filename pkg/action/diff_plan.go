package action

import (
	"os"

	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
)

type DiffLocalPlan struct {
	diff     *Diff
	plandir1 string
	plandir2 string
}

func (d *DiffLocalPlan) Run() error {
	if d.plandir1 == d.plandir2 {
		log.Warn(plan.ErrPlansAreTheSame)
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

	plan1.DiffPlan(plan2, d.diff.ShowSecret, d.diff.Wide)

	return nil
}
