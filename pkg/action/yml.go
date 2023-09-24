package action

import (
	"context"
	"io/fs"

	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var _ Action = (*Yml)(nil)

// Yml is a struct for running 'yml' command.
type Yml struct {
	srcFS     fs.StatFS
	destFS    fs.SubFS
	templater string
}

// Run is the main function for 'yml' command.
func (i *Yml) Run(ctx context.Context) error {
	err := template.Tpl2yml(i.srcFS, i.destFS, "", "", nil, i.templater)
	if err != nil {
		return err
	}

	log.WithField(
		"build plan with next command",
		"helmwave build",
	).Info("📄 YML is ready!")

	return nil
}

// Cmd returns 'yml' *cli.Command.
func (i *Yml) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "yml",
		Usage:  "📄 render helmwave.yml.tpl -> helmwave.yml",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
	}
}

func (i *Yml) flags() []cli.Flag {
	return []cli.Flag{
		flagTplFile(&i.srcFS),
		flagYmlFile(&i.destFS),
		flagTemplateEngine(&i.templater),
	}
}
