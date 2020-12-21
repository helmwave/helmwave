package helmwave

import (
	log "github.com/sirupsen/logrus"
)

func (c *Config) CheckVersion(version string) error {
	if version != c.Version {
		log.Warn("⚠️ Unsupported version ", version)
		log.Debug("🌊 HelmWave version ", c.Version)
	}

	return nil
}
