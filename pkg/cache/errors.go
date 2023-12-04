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

func NewNotCreatedError(dir string, err error) error {
	return &NotCreatedError{Dir: dir, Err: err}
}

func (err NotCreatedError) Error() string {
	return fmt.Sprintf("failed to create cache directory %q: %s", err.Dir, err.Err)
}

func (err NotCreatedError) Unwrap() error {
	return err.Err
}
