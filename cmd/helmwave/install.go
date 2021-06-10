package main

import (
	"github.com/helmwave/helmwave/pkg/action"
	"github.com/urfave/cli/v2"
)

var aInstall = &action.Install{}

var install = &cli.Command{
	Name:    "install",
	Aliases: []string{"apply", "sync", "deploy"},
	Usage:   "🛥 Deploy!",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "plandir",
			Value:       ".helmwave/",
			Usage:       "Path to plandir",
			EnvVars:     []string{"HELMWAVE_PLANDIR"},
			Destination: &aInstall.Plandir,
		},
	},
	Action: toCtx(aInstall.Run),
}
