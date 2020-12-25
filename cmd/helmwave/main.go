package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	helmwave "github.com/zhilyaev/helmwave/pkg/cli"
	"os"
)

func main() {
	app = helmwave.New()
	c := cli.NewApp()
	c.EnableBashCompletion = true
	c.Before = before
	c.CommandNotFound = command404

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

func command404(c *cli.Context, s string) {
	log.Errorf("👻 Command %q not found \n", s)
	os.Exit(127)
}

func before(c *cli.Context) error {
	err := app.InitLogger()
	if err != err {
		return err
	}

	app.InitPlan()
	return nil
}
