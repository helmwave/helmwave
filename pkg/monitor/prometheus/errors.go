package prometheus

import (
	"errors"
	"fmt"
)

var (
	ErrURLEmpty           = errors.New("URL cannot be empty")
	ErrExprEmpty          = errors.New("expression cannot be empty")
	ErrInvalidSuccessMode = errors.New("invalid success mode")

	ErrInvalidResult = errors.New("failed to decode result")

	ErrResultEmpty    = errors.New("result is empty")
	ErrResultNotEmpty = errors.New("result is not empty")
)

type ClientError struct {
	Err error
}

func NewPrometheusClientError(err error) error {
	return &ClientError{Err: err}
}

func (err ClientError) Error() string {
	return fmt.Sprintf("failed to create prometheus client: %s", err.Err)
}

func (err ClientError) Unwrap() error {
	return err.Err
}
