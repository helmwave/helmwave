package action

import (
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/urfave/cli/v2"
)

type Yml struct {
	from, to string
}

func (i *Yml) Run() error {
	return template.Tpl2yml(i.from, i.to, nil)
}

func (i *Yml) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "yml",
		Usage: "📄 Render helmwave.yml.tpl -> helmwave.yml",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "tpl",
				Value:       "helmwave.yml.tpl",
				Usage:       "Main tpl file",
				EnvVars:     []string{"HELMWAVE_TPL"},
				Destination: &i.from,
			},
			flagFile(&i.to),
		},
		Action: toCtx(i.Run),
	}
}
