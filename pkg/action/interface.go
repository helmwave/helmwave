package action

import "github.com/urfave/cli/v2"

type action interface {
	Run() error
	Cmd() *cli.Command
}

func toCtx(a func() error) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return a()
	}
}
