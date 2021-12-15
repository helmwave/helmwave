package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
)

type Rollback struct {
	plandir string
}

func (i *Rollback) Run() error {
	p := plan.New(i.plandir)
	if err := p.Import(); err != nil {
		return err
	}
	return p.Rollback()
}
