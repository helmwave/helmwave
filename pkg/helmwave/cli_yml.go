package helmwave

import (
	"github.com/urfave/cli/v2"
)

func (c *Config) CliYml(ctx *cli.Context) error {
	return c.Tpl.Render()
}
