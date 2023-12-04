package template

import (
	"fmt"
)

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
