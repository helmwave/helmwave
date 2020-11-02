package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	helmwave "github.com/zhilyaev/helmwave/pkg/cli"
	"log"
	"os"
)

func main() {
	app = helmwave.New()
	c := &cli.App{
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
	fmt.Printf("👻 Command '%v' not found \n", s)
}
