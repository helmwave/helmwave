package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

// List is struct for running 'list' command.
type List struct {
	plandir string
}

// Run is main function for 'list' command.
func (l *List) Run() error {
	p := plan.New(l.plandir)
	if err := p.Import(); err != nil {
		return err
	}

	return p.List()
}

// Cmd returns 'list' *cli.Command.
func (l *List) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "👀 List of deployed releases",
		Flags:   l.flags(),
		Action:  toCtx(l.Run),
	}
}

func (l *List) flags() []cli.Flag {
	return []cli.Flag{
		flagPlandir(&l.plandir),
	}
}
