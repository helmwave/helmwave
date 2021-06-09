package main

import (
	"github.com/helmwave/helmwave/pkg/feature"
	"github.com/urfave/cli/v2"
)


var overPlan  = &cli.BoolFlag{
	Name:        "OverPlan",
	Value:       false,
	Usage:       "Allows override plan",
	EnvVars:     []string{"HELMWAVE_PLAN_OVER"},
	Destination: &feature.OverPlan,
}
