package action

import (
	"context"
	"fmt"

	"github.com/helmwave/helmwave/pkg/clictx"
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var _ Action = (*Yml)(nil)

// Yml is a struct for running 'yml' command.
type Yml struct {
	tpl, file string
	templater string
}

// Run is the main function for 'yml' command.
func (i *Yml) Run(ctx context.Context) error {
	cliCtx := clictx.GetCLIFromContext(ctx)
	if cliCtx == nil {
		return fmt.Errorf("failed to get CLI context")
	}

	if i.templater == "" {
		i.templater = cliCtx.String("templater")
	}

	err := template.Tpl2yml(ctx, i.tpl, i.file, nil, i.templater)
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
		Name: "yml",
		// Category: Step0,
		Usage:  "📄 render helmwave.yml.tpl -> helmwave.yml",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
	}
}

func (i *Yml) flags() []cli.Flag {
	return []cli.Flag{
		flagTplFile(&i.tpl),
		flagYmlFile(&i.file),
		flagTemplateEngine(),
	}
}
