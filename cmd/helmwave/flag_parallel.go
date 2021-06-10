package main

import (
	"github.com/urfave/cli/v2"
)

var parallel = &cli.BoolFlag{
	Name:        "parallel",
	Usage:       "It allows parallel mode",
	Value:       true,
	EnvVars:     []string{"HELMWAVE_PARALLEL"},
	Destination: &app.Features.Parallel,
}
