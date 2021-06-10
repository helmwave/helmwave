package action

import (
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/urfave/cli/v2"
)

type Yml struct {
	From string
	To   string
}

func (a *Yml) Run(c *cli.Context) error {
	return template.Tpl2yml(a.From, a.To, nil)
}
