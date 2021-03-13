package main

import (
	helmwaveCli "github.com/helmwave/helmwave/pkg/cli"
	"github.com/helmwave/helmwave/pkg/helmwave"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app = helmwaveCli.New()
	c := cli.NewApp()
	c.EnableBashCompletion = true
	c.CommandNotFound = helmwave.Command404

	c.Usage = "composer for helm"
	c.Version = app.Version
	c.Flags = flags(app)
	c.Commands = commands()
	c.Description = "🏖 This tool helps you compose your helm releases!"

	err := c.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
