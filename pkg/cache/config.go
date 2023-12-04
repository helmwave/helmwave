package cache

import (
	"github.com/adrg/xdg"
	"github.com/urfave/cli/v2"
)

var DefaultConfig = Config{}

type Config struct {
	Home string
}

func (c *Config) Flags() []cli.Flag {
	defaultCache, err := xdg.CacheFile("helmwave")
	if err != nil {
		defaultCache = ""
	}

	return []cli.Flag{
		&cli.PathFlag{
			Name:        "cache-dir",
			Usage:       "base directory for cache",
			Value:       defaultCache,
			EnvVars:     []string{"HELMWAVE_CACHE_DIR", "HELMWAVE_CACHE_HOME"},
			Destination: &c.Home,
		},
	}
}

// Run initializes cache.
func (c *Config) Run(_ *cli.Context) error {
	return c.Init()
}

// Init initializes cache.
func (c *Config) Init() error {
	return nil
}
