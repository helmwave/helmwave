package action

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Action is an interface for all actions.
type Action interface {
	Run(context.Context) error
	Cmd() *cli.Command
	flags() []cli.Flag
}

// toCtx is a wrapper for urfave v2.
func toCtx(a func(context.Context) error) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		ctx := getContextWithFlags(c)

		return a(ctx)
	}
}

func getContextWithFlags(c *cli.Context) context.Context {
	ctx := c.Context
	for _, flagName := range c.FlagNames() {
		g := c.Value(flagName)
		log.WithField("name", flagName).WithField("value", g).Trace("adding flag to action context.Context")
		ctx = context.WithValue(ctx, flagName, g) //nolint:staticcheck // weird issue, we won't have any collisions with strings
	}

	//nolint:staticcheck // same
	ctx = context.WithValue(ctx, "cli", c)

	return ctx
}
