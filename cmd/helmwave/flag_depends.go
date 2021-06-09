package main

import (
	"github.com/helmwave/helmwave/pkg/feature"
	"github.com/urfave/cli/v2"
)


var depends  = &cli.BoolFlag{
	Name:        "enable-dependencies",
	Usage:       "Enable dependencies",
	Value:       false,
	EnvVars:     []string{"HELMWAVE_ENABLE_DEPENDENCIES"},
	Destination: &feature.Dependencies,
}
