package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/helmwave/helmwave/pkg/action"
	"github.com/helmwave/helmwave/pkg/cache"
	logSetup "github.com/helmwave/helmwave/pkg/log"
	helmwave "github.com/helmwave/helmwave/pkg/version"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

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

func main() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	c := CreateApp()

	defer recoverPanic()

	if err := c.Run(os.Args); err != nil {
		log.Fatal(err) //nolint:gocritic // we try to recover panics, not regular command errors
	}
}

func recoverPanic() {
	if r := recover(); r != nil {
		switch r.(type) {
		case CommandNotFoundError:
			log.Error(r)
			log.Exit(127)
		default:
			log.Panic(r)
		}
	}
}

// CreateApp creates *cli.App with all commands.
func CreateApp() *cli.App {
	c := cli.NewApp()

	c.EnableBashCompletion = true
	c.Usage = "true release management for helm"
	c.Version = helmwave.Version
	c.Description = "This tool helps you compose your helm releases!\n" +
		"0. $ helmwave yml\n" +
		"1. $ helmwave build\n" +
		"2. $ helmwave up\n"

	c.Before = func(ctx *cli.Context) error {
		if ctx.Bool("cancel-on-kill") {
			// cancel the ctx if it gets these signals
			ctx.Context = cancelCtxOnSignal(ctx.Context, syscall.SIGTERM, syscall.SIGINT)
		}
		err := logSetup.Default.Run(ctx)
		if err != nil {
			return err
		}

		err = cache.DefaultConfig.Run(ctx)
		if err != nil {
			return err
		}

		return nil
	}
	c.Flags = append(logSetup.Default.Flags(), cache.DefaultConfig.Flags()...)
	c.Flags = append(c.Flags, cancelFlag())

	c.Commands = commands
	c.CommandNotFound = command404

	return c
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

func cancelFlag() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  "cancel-on-kill",
		Usage: "cancel helm on SigINT,SigTERM",

		Value:   false,
		EnvVars: []string{"HELMWAVE_HANDLE_SIGNAL"},
		Hidden:  true, // experimental
	}
}

type signalCtx struct {
	context.Context
}

func (ctx signalCtx) Err() error {
	// Unwrap for readablity, if ctx cancelled by reciving a signal
	if sigErr, ok := context.Cause(ctx.Context).(signalErr); ok {
		return sigErr
	}
	return ctx.Context.Err()
}

type signalErr struct {
	sig os.Signal
}

func (e signalErr) Error() string {
	return "context canceled: got signal " + e.sig.String()
}

// cancelCtxOnSignal closes Done channel when one of the listed signals arrives
func cancelCtxOnSignal(parent context.Context, signals ...os.Signal) (ctx context.Context) {
	ctx, cancel := context.WithCancelCause(parent)
	ctx = signalCtx{ctx} // wrap for err handling

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	if ctx.Err() == nil {
		go func() {
			select {
			case sig := <-ch:
				// log.Warn("got signal " + sig.String())
				cancel(signalErr{sig})
			case <-ctx.Done():
			}
		}()
	}
	return ctx
}
