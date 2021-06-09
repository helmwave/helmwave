package main

import (
	helmwaveCli "github.com/helmwave/helmwave/pkg/cli"
	"github.com/helmwave/helmwave/pkg/helmwave"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

var app *helmwave.Config

var commands = []*cli.Command{
	version,
	yml,
	install,
	uninstall,
	status,
	list,
	plan,
	manifest,
}

func main() {
	app = helmwaveCli.New()
	c := cli.NewApp()
	c.EnableBashCompletion = true
	c.CommandNotFound = helmwave.Command404
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
