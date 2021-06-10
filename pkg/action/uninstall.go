package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
)

type Uninstall struct {
	Plandir string
}

func (i *Uninstall) Run() error {
	p := plan.New(i.Plandir)
	return p.Apply()
}
