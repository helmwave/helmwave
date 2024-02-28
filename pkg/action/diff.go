package action

import (
	"github.com/databus23/helm-diff/v3/diff"
	logSetup "github.com/helmwave/helmwave/pkg/log"
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
	*diff.Options
	kindSuppressHelper cli.StringSlice
	findRenamesHelper  float64
}

func (d *Diff) FixFields() {
	d.OutputFormat = logSetup.Default.Format()
	d.SuppressedKinds = d.kindSuppressHelper.Value()
	d.FindRenames = float32(d.findRenamesHelper)
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
		Before: func(q *cli.Context) error {
			d.FixFields()

			return nil
		},
		Subcommands: []*cli.Command{
			plan.Cmd(),
			live.Cmd(),
		},
	}
}

// flags return flag set of CLI urfave.
func (d *Diff) flags() []cli.Flag {
	d.Options = &diff.Options{}

	self := []cli.Flag{
		flagDiffWide(&d.OutputContext),
		flagDiffShowSecret(&d.ShowSecrets),
		&cli.BoolFlag{
			Name:        "strip-trailing-cr",
			Usage:       "strip trailing carriage return on input",
			Value:       false,
			EnvVars:     []string{"HELMWAVE_DIFF_STRIP_TRAILING_CR"},
			Destination: &d.StripTrailingCR,
		},
		&cli.Float64Flag{
			Name:        "find-renames",
			Usage:       "enable rename detection if set to any value greater than 0.",
			Value:       0,
			EnvVars:     []string{"HELMWAVE_DIFF_FIND_RENAMES"},
			Destination: &d.findRenamesHelper,
		},
		&cli.StringSliceFlag{
			Name:  "suppress",
			Usage: "allows suppression of the values listed in the diff output (\"Secret\" for example)",
			// Value: cli.NewStringSlice("Secret"),
			EnvVars:     []string{"HELMWAVE_DIFF_SUPPRESS"},
			Destination: &d.kindSuppressHelper,
		},
	}

	return self
}
