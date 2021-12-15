package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
)

type Validate struct {
	plandir string
}

func (l *Validate) Run() error {
	p := plan.New(l.plandir)
	if err := p.Import(); err != nil {
		return err
	}
	return p.ValidateValues()
}
