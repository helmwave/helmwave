package main

import (
	"github.com/helmwave/helmwave/pkg/action"
	"github.com/urfave/cli/v2"
)

var aDiff = &action.Diff{}

var diff = &cli.Command{
	Name:    "diff",
	Usage:   "📜 Diff 2 plans",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "plan B",
			Value:       ".helmwave/",
			Usage:       "Path to plandir A",
			EnvVars:     []string{"HELMWAVE_PLANDIR_A", "HELMWAVE_PLANDIR"},
			Destination: &aDiff.Plandir1,
		},
		&cli.StringFlag{
			Name:        "plan A",
			Usage:       "Path to plandir B",
			EnvVars:     []string{"HELMWAVE_PLANDIR_B"},
			Destination: &aDiff.Plandir1,
		},
	},
	Action: aDiff.Run,
}
