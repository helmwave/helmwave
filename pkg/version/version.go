package version

import (
	log "github.com/sirupsen/logrus"
)

const (
	// Version is helmwave binary version.
	// It will override by goreleaser during release.
	Version = "dev"
)

// Check compares helmwave versions and logs difference.
func Check(a, b string) {
	if a != b {
		log.Warn("⚠️ Unsupported version ", b)
		log.Debug("🌊 HelmWave version ", a)
	}
}
