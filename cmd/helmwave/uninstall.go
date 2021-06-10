package main

import (
	"github.com/helmwave/helmwave/pkg/action"
	"github.com/urfave/cli/v2"
)

var aUninstall = &action.Uninstall{}
var uninstall = &cli.Command	{
	Name:    "uninstall",
	Aliases: []string{"destroy", "delete", "del", "rm", "remove"},
	Usage:   "Delete deployed releases",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "plandir",
			Value:       ".helmwave/",
			Usage:       "Path to plandir",
			EnvVars:     []string{"HELMWAVE_PLANDIR"},
			Destination: &aUninstall.Plandir,
		},
	},
	Action: toCtx(aUninstall.Run),
}
