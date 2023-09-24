package template

import (
	"errors"
	"fmt"
)

var ErrInvalidFilesystem = errors.New("filesystem not supported")

type SOPSDecodeError struct {
	Err error
}

func NewSOPSDecodeError(err error) error {
	return &SOPSDecodeError{Err: err}
}

func (err SOPSDecodeError) Error() string {
	return fmt.Sprintf("failed to decode values file with SOPS: %s", err.Err)
}

func (err SOPSDecodeError) Unwrap() error {
	return err.Err
}

func (SOPSDecodeError) Is(target error) bool {
	switch target.(type) {
	case SOPSDecodeError, *SOPSDecodeError:
		return true
	default:
		return false
	}
}
