package action

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

var _ Action = (*GenSchema)(nil)

// GenSchema is a struct for running 'GenSchema' command.
type GenSchema struct{}

// Run is the main function for 'GenSchema' command.
func (i *GenSchema) Run(ctx context.Context) error {
	s, err := json.Marshal(plan.GenSchema())
	if err != nil {
		return err
	}

	fmt.Println(string(s)) //nolint:forbidigo

	return nil
}

// Cmd returns 'GenSchema' *cli.Command.
func (i *GenSchema) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "schema",
		Usage:  "generate json schema",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
	}
}

func (i *GenSchema) flags() []cli.Flag {
	return nil
}
