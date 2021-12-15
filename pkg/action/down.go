package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
)

type Down struct {
	plandir string
}

func (i *Down) Run() error {
	p := plan.New(i.plandir)
	if err := p.Import(); err != nil {
		return err
	}
	return p.Destroy()
}
