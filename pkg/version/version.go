package version

import (
	log "github.com/sirupsen/logrus"
)

// Version is helmwave binary version.
// It should be var not const.
// It will override by goreleaser during release.
// -X github.com/helmwave/helmwave/pkg/version.Version={{ .Version }}.
//
// nolintlint:gochecknoglobals
var Version = "dev" //nolint:gochecknoglobals

// Check compares helmwave versions and logs difference.
func Check(a, b string) {
	if a != b {
		log.Warn("⚠️ Unsupported version ", b)
		log.Debug("🌊 HelmWave version ", a)
	}
}
