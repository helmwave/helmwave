package version

import (
	log "github.com/sirupsen/logrus"
)

func Check(a, b string) {
	if a != b {
		log.Warn("⚠️ Unsupported version ", b)
		log.Debug("🌊 HelmWave version ", a)
	}
}

var Version = "1.0.0+dev"
