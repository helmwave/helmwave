package fileref

import (
	"errors"

	"helm.sh/helm/v3/pkg/storage/driver"
)

var (
	ErrNameEmpty = errors.New("release name is empty")

	// ErrPendingRelease is an error for fail strategy that release is in pending status.
	ErrPendingRelease = errors.New("release is in pending status")

	// ErrValuesNotExist is returned when values can't be used and are skipped.
	ErrValuesNotExist = errors.New("values file doesn't exist")

	// ErrNotFound is an error for not found release.
	ErrNotFound = driver.ErrReleaseNotFound

	// ErrFoundMultiple is an error for multiple releases found by name.
	ErrFoundMultiple = errors.New("found multiple releases o_0")

	// ErrDepFailed is an error thrown when dependency release fails.
	ErrDepFailed = errors.New("dependency failed")

	ErrUnknownFormat = errors.New("unknown format")

	ErrDigestNotMatch = errors.New("chart digest doesn't match")
)
