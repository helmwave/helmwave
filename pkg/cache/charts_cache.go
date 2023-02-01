package cache

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
)

//nolint:gochecknoglobals
var ChartsCache = CacheConfig{}

type CacheConfig struct {
	cacheDir string
}

func (c *CacheConfig) Init(dir string) {
	c.cacheDir = dir
}

func (c *CacheConfig) IsEnabled() bool {
	return c.cacheDir != ""
}

func (c *CacheConfig) FindInCache(chart string, version string) (string, error) {
	if !c.IsEnabled() {
		return "", fmt.Errorf("cache is disabled")
	}

	chartName := filepath.Base(chart)
	chartFile := fmt.Sprintf("%s/%s-%s.tgz", c.cacheDir, chartName, version)

	_, err := os.Stat(chartFile)
	if err == nil {
		return chartFile, nil
	}

	return "", fmt.Errorf("chart not found")
}

func (c *CacheConfig) AddToCache(file string) {
	if !c.IsEnabled() {
		return
	}

	if err := helper.CopyFile(file, c.cacheDir); err != nil {
		log.Warn(fmt.Errorf("failed to cache chart %s: %w", file, err))
	}
}
