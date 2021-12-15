package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
)

type List struct {
	plandir string
}

func (l *List) Run() error {
	p := plan.New(l.plandir)
	if err := p.Import(); err != nil {
		return err
	}
	return p.List()
}
