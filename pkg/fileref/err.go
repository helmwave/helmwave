package fileref

import (
	"errors"
)

var (
	// ErrValuesNotExist is returned when values can't be used and are skipped.
	ErrValuesNotExist = errors.New("file doesn't exist")

	ErrUnknownFormat = errors.New("unknown format")
)
