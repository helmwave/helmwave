package sops

import (
	"fmt"
)

type DecodeError struct {
	Err error
}

func NewDecodeError(err error) error {
	return &DecodeError{Err: err}
}

func (err DecodeError) Error() string {
	return fmt.Sprintf("failed to decode values file with SOPS: %s", err.Err)
}

func (err DecodeError) Unwrap() error {
	return err.Err
}
