package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/helmwave/helmwave/pkg/helper"

	"github.com/helmwave/helmwave/pkg/action"
	"github.com/helmwave/helmwave/pkg/cache"
	logSetup "github.com/helmwave/helmwave/pkg/log"
	helmwave "github.com/helmwave/helmwave/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	helper.Dotenv()

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
			// https://tldp.org/LDP/abs/html/exitcodes.html
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

	c.Before = before
	c.Flags = action.GlobalFlags()
	c.Commands = commands
	c.CommandNotFound = command404

	return c
}

// cancelCtxOnSignal closes Done channel when one of the listed signals arrives.
func cancelCtxOnSignal(parent context.Context, signals ...os.Signal) (ctx context.Context) {
	ctx, cancel := context.WithCancelCause(parent) //nolint:govet

	context.AfterFunc(ctx, func() {
		if err := context.Cause(ctx); err != nil {
			log.Error(err)
		}
	})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	if ctx.Err() == nil {
		go func() {
			select {
			case sig := <-ch:
				cancel(fmt.Errorf("got signal %v", sig))

			case <-ctx.Done():
			}
		}()
	}

	return ctx //nolint:govet // lostcancel: this return statement may be reached without using the cancel var defined on line 67
}

func before(ctx *cli.Context) error {
	if ctx.Bool("handle-signal") {
		ctx.Context = cancelCtxOnSignal(ctx.Context, syscall.SIGTERM, syscall.SIGINT)
	}

	// Init flags first
	err := logSetup.Default.Run(ctx)
	if err != nil {
		return err
	}
	err = cache.Default.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}
