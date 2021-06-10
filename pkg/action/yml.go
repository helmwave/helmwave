package action

import (
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/urfave/cli/v2"
)

type Yml struct {
	From string
	To   string
}

func (i *Yml) Run() error {
	return template.Tpl2yml(i.From, i.To, nil)
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
				EnvVars:     []string{"HELMWAVE_TPL_FILE"},
				Destination: &i.From,
			},
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Value:       "helmwave.yml",
				Usage:       "Main yml file",
				EnvVars:     []string{"HELMWAVE_FILE", "HELMWAVE_YAML_FILE", "HELMWAVE_YML_FILE"},
				Destination: &i.To,
			},
		},
		Action: toCtx(i.Run),
	}
}
