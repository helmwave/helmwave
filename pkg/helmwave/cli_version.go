package helmwave

import "github.com/urfave/cli/v2"

func (c *Config) CliVersion(ctx *cli.Context) error {
	cli.ShowVersion(ctx)
	return nil
}
