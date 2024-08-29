package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"path/filepath"

	"github.com/adrg/xdg"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var Default = Config{}

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

func (s *Config) GetRemoteSourcePath(remoteSource *url.URL) string {
	u := *remoteSource
	u.RawQuery = ""

	hasher := sha256.New()
	hasher.Write([]byte(u.String()))
	hash := hex.EncodeToString(hasher.Sum(nil))

	p := filepath.Join(s.Home, "remote-source", hash)

	log.WithField("remote source", remoteSource.String()).
		WithField("cache path", p).
		Info("using cache for remote source")

	return p
}
