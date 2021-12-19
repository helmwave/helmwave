package action

import "github.com/urfave/cli/v2"

//nolint:unused
type Action interface {
	Run() error
	Cmd() *cli.Command
}

// toCtx wrapper of urfave v2
func toCtx(a func() error) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return a()
	}
}
