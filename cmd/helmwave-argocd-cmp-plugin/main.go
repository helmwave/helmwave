package main

import (
	"os"

	"github.com/helmwave/helmwave/pkg/action"
	"github.com/helmwave/helmwave/pkg/cache"
	"github.com/helmwave/helmwave/pkg/helper"
	logSetup "github.com/helmwave/helmwave/pkg/log"
	helmwave "github.com/helmwave/helmwave/pkg/version"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	c := CreateApp()

	if err := c.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// CreateApp creates *cli.App with all commands.
func CreateApp() *cli.App {
	c := cli.NewApp()

	c.Usage = "true release management for helm"
	c.Version = helmwave.Version

	c.Before = func(ctx *cli.Context) error {
		logSetup.Default.ArgoCDMode()

		err := logSetup.Default.Run(ctx)
		if err != nil {
			return err
		}

		err = cache.DefaultConfig.Run(ctx)
		if err != nil {
			return err
		}

		helper.Helm.RegistryConfig = ".helm-config/registry/config.json"
		helper.Helm.RepositoryConfig = ".helm-config/repositories.yaml"
		helper.Helm.RepositoryCache = ".helm-cache/repository"

		return nil
	}

	act := &action.Build{}
	act.ArgoCDMode()

	// we just need the manifests, no other actions required
	c.Action = act.Cmd().Action

	return c
}
