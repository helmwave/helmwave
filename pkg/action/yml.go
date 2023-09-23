package action

import (
	"context"
	"io/fs"
	"net/url"
	"os"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var _ Action = (*Yml)(nil)

// Yml is a struct for running 'yml' command.
type Yml struct {
	srcFS     fs.FS
	destFS    fsimpl.WriteableFS
	tpl, file string
	templater string
}

func getBaseFS() fsimpl.WriteableFS {
	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})

	return baseFS.(fsimpl.WriteableFS) //nolint:forcetypeassert
}

// Run is the main function for 'yml' command.
func (i *Yml) Run(ctx context.Context) error {
	// TODO: get filesystems dynamically from args
	i.srcFS = getBaseFS()
	i.destFS = getBaseFS()

	err := template.Tpl2yml(i.srcFS, i.destFS, i.tpl, i.file, nil, i.templater)
	if err != nil {
		return err
	}

	log.WithField(
		"build plan with next command",
		"helmwave build -f "+i.file,
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
		flagTplFile(&i.tpl),
		flagYmlFile(&i.file),
		flagTemplateEngine(&i.templater),
	}
}
