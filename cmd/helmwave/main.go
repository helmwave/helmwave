package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	helmwave "github.com/zhilyaev/helmwave/pkg/cli"
	"os"
)

func main() {
	app = helmwave.New()
	c := &cli.App{
		Before:          before,
		CommandNotFound: command404,
		Name:            "🌊 HelmWave",
		Usage:           "composer for helm",
		Version:         app.Version,
		Authors:         authors(),
		Flags:           flags(app),
		Commands:        commands(),
		Description:     "🏖 This tool helps you compose your helm releases!",
	}

	err := c.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func command404(c *cli.Context, s string) {
	log.Errorf("👻 Command '%v' not found \n", s)
	os.Exit(127)
}

func before(c *cli.Context) error {
	return app.InitLogger()
}
