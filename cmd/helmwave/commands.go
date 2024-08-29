package main

import (
	"fmt"

	"github.com/helmwave/helmwave/pkg/action"
	helmwave "github.com/helmwave/helmwave/pkg/version"
	"github.com/urfave/cli/v2"
)

// commands is a registration list for commands.
var commands = []*cli.Command{
	new(action.Build).Cmd(),
	new(action.Diff).Cmd(),
	new(action.Up).Cmd(),
	new(action.List).Cmd(),
	new(action.Rollback).Cmd(),
	new(action.Status).Cmd(),
	new(action.Down).Cmd(),
	new(action.Validate).Cmd(),
	new(action.Yml).Cmd(),
	new(action.GenSchema).Cmd(),
	new(action.Graph).Cmd(),
	version(),
	completion(),
}

func version() *cli.Command {
	return &cli.Command{
		Name:     "version",
		Aliases:  []string{"ver"},
		Category: action.Step_,
		Usage:    "show shorts version",
		Action: func(c *cli.Context) error {
			fmt.Println(helmwave.Version) //nolint:forbidigo // we need to use fmt.Println here

			return nil
		},
	}
}

// CommandNotFoundError is return when CLI command is not found.
type CommandNotFoundError struct {
	Command string
}

func (e CommandNotFoundError) Error() string {
	return fmt.Sprintf("👻 Command %q not found", e.Command)
}

func command404(_ *cli.Context, s string) {
	err := CommandNotFoundError{
		Command: s,
	}
	panic(err)
}
