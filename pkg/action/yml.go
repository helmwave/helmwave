package action

import (
	"github.com/helmwave/helmwave/pkg/template"
)

type Yml struct {
	From string
	To   string
}

func (a *Yml) Run() error {
	return template.Tpl2yml(a.From, a.To, nil)
}
