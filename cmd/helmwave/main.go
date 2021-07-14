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
	new(action.Up).Cmd(),
	new(action.List).Cmd(),
	new(action.Rollback).Cmd(),
	new(action.Down).Cmd(),
	new(action.Validate).Cmd(),
	new(action.Yml).Cmd(),
	completion(),
}

func main() {
	c := cli.NewApp()
	c.EnableBashCompletion = true
	c.Usage = "is like docker-compose for helm"
	c.Version = helmwave.Version
	c.Description =
		"This tool helps you compose your helm releases!\n" +
			"0. $ helmwave yml\n" +
			"1. $ helmwave build\n" +
			"2. $ helmwave apply\n"

	logSet := logSetup.Settings{}
	c.Before = logSet.Run
	c.Flags = logSet.Flags()

	c.Commands = commands

	err := c.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
