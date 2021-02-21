package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/zhilyaev/helmwave/pkg/helmwave"
	helmwaveCli "github.com/zhilyaev/helmwave/pkg/cli"
	"os"
)

func main() {
	app = helmwaveCli.New()
	c := cli.NewApp()
	c.EnableBashCompletion = true
	c.CommandNotFound = helmwave.Command404

	c.Usage = "composer for helm"
	c.Version = app.Version
	c.Authors = authors()
	c.Flags = flags(app)
	c.Commands = commands()
	c.Description = "🏖 This tool helps you compose your helm releases!"

	err := c.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

