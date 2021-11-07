package action

import (
	"github.com/urfave/cli/v2"
)

type Diff struct {
	Wide       int
	ShowSecret bool
}

func (d *Diff) Cmd() *cli.Command {

	plan := DiffPlans{diff: d}
	live := DiffLive{diff: d}

	return &cli.Command{
		Name:    "diff",
		Usage:   "🆚 Show Differences",
		Aliases: []string{"vs"},
		Flags:   d.flags(),
		Subcommands: []*cli.Command{
			plan.Cmd(),
			live.Cmd(),
		},
	}
}

func (d *Diff) flags() []cli.Flag {
	return []cli.Flag{
		flagDiffWide(&d.Wide),
		flagDiffShowSecret(&d.ShowSecret),
	}
}
