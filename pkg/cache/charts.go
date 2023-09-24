package cache

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
)

//nolint:gochecknoglobals
var ChartsCache = Config{}

type Config struct {
	cacheFS  fsimpl.WriteableFS
	cacheDir string
	lock     sync.RWMutex
}

func (c *Config) Init(cacheFS fs.FS, cacheDir string) error {
	if cacheFS == nil || cacheDir == "" {
		return nil
	}

	writeableCacheFS, ok := cacheFS.(fsimpl.WriteableFS)
	if !ok {
		return ErrNotWriteableFS
	}

	c.cacheFS = writeableCacheFS
	c.cacheDir = cacheDir

	if !helper.IsExists(c.cacheFS, c.cacheDir) {
		if err := c.cacheFS.MkdirAll(c.cacheDir, 0o755); err != nil {
			return NewNotCreatedError(c.cacheDir, err)
		}
	}

	return nil
}

func (c *Config) IsEnabled() bool {
	return c.cacheFS != nil && c.cacheDir != ""
}

func (c *Config) FindInCache(chart, version string) (string, error) {
	if !c.IsEnabled() {
		return "", ErrCacheDisabled
	}

	chartName := filepath.Base(chart)
	chartFile := helper.FilepathJoin(c.cacheDir, fmt.Sprintf("%s-%s.tgz", chartName, version))

	c.lock.RLock()
	defer c.lock.RUnlock()

	_, err := c.cacheFS.Stat(chartFile)
	if err == nil {
		return chartFile, nil
	}

	return "", ErrChartNotFound
}

func (c *Config) AddToCache(f fs.FS, file string) {
	if !c.IsEnabled() {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	if err := helper.CopyFile(f, c.cacheFS, file, c.cacheDir); err != nil {
		log.WithError(err).Warnf("failed to cache chart %s", file)
	}
}
