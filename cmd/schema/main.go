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
	c.Commands = commands
	c.Version = helmwave.Version

	if err := c.Run(os.Args); err != nil {
		log.Fatal(err) //nolint:gocritic // we try to recover panics, not regular command errors
	}
}

var commands = []*cli.Command{
	new(action.GenSchema).Cmd(),
}
