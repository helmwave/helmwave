package action

import (
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type Yml struct {
	tpl, file string
	templater string
}

func (i *Yml) Run() error {
	err := template.Tpl2yml(i.tpl, i.file, nil, i.templater)
	if err != nil {
		return err
	}

	log.WithField(
		"build plan with next command",
		"helmwave build -f "+i.file,
	).Info("📄 YML is ready!")

	return nil
}

func (i *Yml) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "yml",
		Usage:  "📄 Render helmwave.yml.tpl -> helmwave.yml",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
	}
}

func (i *Yml) flags() []cli.Flag {
	return []cli.Flag{
		flagTplFile(&i.tpl),
		flagYmlFile(&i.file),
		flagTemplateEngine(&i.templater),
	}
}
