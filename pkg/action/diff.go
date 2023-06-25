package action

import (
	"github.com/urfave/cli/v2"
)

const (
	// DiffModeLive is a subcommand name for diffing manifests in plan with actually running manifests in k8s.
	DiffModeLive = "live"

	// DiffModeLocal is a subcommand name for diffing manifests in two plans.
	DiffModeLocal = "local"

	// DiffModeNone is a subcommand name for skipping diffing.
	DiffModeNone = "none"
)

// Diff is a struct for running 'diff' commands.
type Diff struct {
	ThreeWayMerge bool
	ShowSecret    bool
	Wide          int
}

// Cmd returns 'diff' *cli.Command.
func (d *Diff) Cmd() *cli.Command {
	plan := DiffLocal{diff: d}
	live := DiffLive{diff: d}

	return &cli.Command{
		Name:    "diff",
		Usage:   "🆚 show differences",
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
