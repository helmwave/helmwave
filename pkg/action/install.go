package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
)

type Install struct {
	Plandir string
}

func (i *Install) Run() error {
	p := plan.New(i.Plandir)
	return p.Apply()
}
