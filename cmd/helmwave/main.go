package main

import (
	"os"

	"github.com/helmwave/helmwave/pkg/action"
	logSetup "github.com/helmwave/helmwave/pkg/log"
	helmwave "github.com/helmwave/helmwave/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var commands = []*cli.Command{
	new(action.Build).Cmd(),
	new(action.Diff).Cmd(),
	new(action.Install).Cmd(),
	new(action.List).Cmd(),
	new(action.Rollback).Cmd(),
	new(action.Uninstall).Cmd(),
	new(action.Validate).Cmd(),
	new(action.Yml).Cmd(),
}

func main() {
	c := cli.NewApp()
	c.EnableBashCompletion = true
	c.Usage = "composer for helm"
	c.Version = helmwave.Version
	c.Description = "This tool helps you compose your helm releases!"

	logSet := logSetup.Settings{}
	c.Before = logSet.Run
	c.Flags = logSet.Flags()

	c.Commands = commands

	err := c.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
