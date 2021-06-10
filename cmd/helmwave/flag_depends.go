package main

import (
	"github.com/urfave/cli/v2"
)

var dependsOn = &cli.BoolFlag{
	Name:        "enable-dependencies",
	Usage:       "Enable dependencies",
	Value:       false,
	EnvVars:     []string{"HELMWAVE_ENABLE_DEPENDENCIES"},
	Destination: &app.Features.DependsOn,
}
