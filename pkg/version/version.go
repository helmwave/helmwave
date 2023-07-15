package version

import (
	log "github.com/sirupsen/logrus"
)

// Version is a helmwave binary version.
// It should be a var not const.
// It will override by goreleaser during release.
// -X github.com/helmwave/helmwave/pkg/version.Version={{ .Version }}.
var Version = "dev"

// validate compares helmwave versions.
func validate(a, b string) bool {
	if a != b {
		log.Warnf("⚠️ yaml version is %s but binary version is %s", a, b)

		return false
	}

	log.Debug("✅ yaml version is equal to binary version")

	return true
}

// Validate compare version.
func Validate(a string) bool {
	return validate(a, Version)
}
