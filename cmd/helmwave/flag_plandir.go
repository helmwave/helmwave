package main

import "github.com/urfave/cli/v2"


var plandir  = &cli.StringFlag{
	Name:        "plandir",
	Value:       ".helmwave/",
	Usage:       "Path to planfile",
	EnvVars:     []string{"HELMWAVE_PLAN_DIR", "HELMWAVE_PLANDIR"},
	Destination: &app.Plandir,
}
