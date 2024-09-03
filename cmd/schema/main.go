package main

import (
	"log"
	"os"

	"github.com/helmwave/helmwave/pkg/action"
	helmwave "github.com/helmwave/helmwave/pkg/version"
	"github.com/urfave/cli/v2"
)

func main() {
	c := cli.NewApp()
	c.Usage = "just generates json schema for helmwave support"
	c.Commands = commands
	c.Version = helmwave.Version

	if err := c.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var commands = []*cli.Command{
	new(action.GenSchema).Cmd(),
}
