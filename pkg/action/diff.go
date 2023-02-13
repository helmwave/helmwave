package action

import (
	"github.com/urfave/cli/v2"
)

// Diff is struct for running 'diff' commands.
type Diff struct {
	ThreeWayMerge bool
	ShowSecret    bool
	Wide          int
}

// Cmd returns 'diff' *cli.Command.
func (d *Diff) Cmd() *cli.Command {
	plan := DiffLocalPlan{diff: d}
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

// flags return flag set of CLI urfave.
func (d *Diff) flags() []cli.Flag {
	return []cli.Flag{
		flagDiffWide(&d.Wide),
		flagDiffShowSecret(&d.ShowSecret),
		flagDiffThreeWayMerge(&d.ThreeWayMerge),
	}
}
