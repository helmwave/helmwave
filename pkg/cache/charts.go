package cache

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
)

//nolintlint:gochecknoglobals
var ChartsCache = Config{}

type Config struct {
	cacheDir string
	lock     sync.RWMutex
}

func (c *Config) Init(dir string) error {
	c.cacheDir = dir
	if dir == "" {
		return nil
	}
	if !helper.IsExists(dir) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create cache directory %s: %w", dir, err)
		}
	}

	return nil
}

func (c *Config) IsEnabled() bool {
	return c.cacheDir != ""
}

func (c *Config) FindInCache(chart string, version string) (string, error) {
	if !c.IsEnabled() {
		return "", fmt.Errorf("cache is disabled")
	}

	chartName := filepath.Base(chart)
	chartFile := path.Join(c.cacheDir, fmt.Sprintf("%s-%s.tgz", chartName, version))

	c.lock.RLock()
	defer c.lock.RUnlock()

	_, err := os.Stat(chartFile)
	if err == nil {
		return chartFile, nil
	}

	return "", fmt.Errorf("chart not found")
}

func (c *Config) AddToCache(file string) {
	if !c.IsEnabled() {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	if err := helper.CopyFile(file, c.cacheDir); err != nil {
		log.Warn(fmt.Errorf("failed to cache chart %s: %w", file, err))
	}
}
