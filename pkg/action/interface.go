package action

import "github.com/urfave/cli/v2"

type action interface {
	Run(c *cli.Context) error
}
