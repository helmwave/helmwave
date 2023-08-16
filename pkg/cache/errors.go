package cache

import (
	"errors"
	"fmt"
)

var (
	ErrCacheDisabled = errors.New("cache is disabled")
	ErrChartNotFound = errors.New("chart not found")
)

type NotCreatedError struct {
	Err error
	Dir string
}

func (err NotCreatedError) Error() string {
	return fmt.Sprintf("failed to create cache directory %s: %v", err.Dir, err.Err)
}

func (err NotCreatedError) Unwrap() error {
	return err.Err
}

//nolint:errorlint
func (NotCreatedError) Is(target error) bool {
	switch target.(type) {
	case NotCreatedError, *NotCreatedError:
		return true
	default:
		return false
	}
}
