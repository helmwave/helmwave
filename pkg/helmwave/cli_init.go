package helmwave

import "github.com/urfave/cli/v2"

func (c *Config) InitApp(ctx *cli.Context) error {
	err := c.InitLogger()
	if err != nil {
		return err
	}

	c.InitPlan()
	return nil
}
