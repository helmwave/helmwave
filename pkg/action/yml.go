package action

import (
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type Yml struct {
	from, to string
}

func (i *Yml) Run() error {
	err := template.Tpl2yml(i.from, i.to, nil)
	if err != nil {
		return err
	}

	log.WithField(
		"build plan with next command",
		"helmwave build",
	).Info("📄 YML is ready!")

	return nil
}

func (i *Yml) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "yml",
		Usage: "📄 Render helmwave.yml.tpl -> helmwave.yml",
		Flags: []cli.Flag{
			flagTplFile(&i.from),
			flagYmlFile(&i.to),
		},
		Action: toCtx(i.Run),
	}
}
