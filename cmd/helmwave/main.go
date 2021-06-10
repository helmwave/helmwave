package main

import (
	"github.com/helmwave/helmwave/pkg/action"
	"github.com/helmwave/helmwave/pkg/helmwave"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

var app *helmwave.Config

var commands = []*cli.Command{
	new(action.Diff).Cmd(),
	new(action.Install).Cmd(),
	new(action.List).Cmd(),
	new(action.Status).Cmd(),
	new(action.Uninstall).Cmd(),
	new(action.Yml).Cmd(),
}


func main() {
	app = helmwave.New()
	c := cli.NewApp()
	c.EnableBashCompletion = true
	c.Usage = "composer for helm"
	c.Version = app.Version
	c.Description = "This tool helps you compose your helm releases!"

	// Default flags and commands
	c.Flags = flagsLog
	c.Commands = commands

	err := c.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
