package action

import (
	"context"

	"github.com/helmwave/helmwave/pkg/clictx"
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
		ctx := clictx.CLIContextToContext(c)

		return a(ctx)
	}
}
