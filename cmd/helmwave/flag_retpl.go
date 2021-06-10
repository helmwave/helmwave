package main

import (
	"github.com/urfave/cli/v2"
)

var reTpl = &cli.BoolFlag{
	Name:        "OverPlan",
	Value:       false,
	Usage:       "Allows re template plan",
	EnvVars:     []string{"HELMWAVE_RETPL"},
	Destination: &app.Features.ReTpl,
}
