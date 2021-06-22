package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

type Manifest struct {
	plandir string
}

func (i *Manifest) Run() error {
	p := plan.New(i.plandir)
	if err := p.Import(); err != nil {
		return err
	}
	return p.Apply()
}

func (i *Manifest) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "manifest",
		Usage: "Makes manifests",
		Flags: []cli.Flag{
			flagPlandir(&i.plandir),
		},
		Action: toCtx(i.Run),
	}
}
