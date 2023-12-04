package prometheus

import (
	"errors"
	"fmt"
)

var (
	ErrURLEmpty = errors.New("URL cannot be empty")

	ErrExprEmpty = errors.New("expression cannot be empty")

	ErrResultNotVector = errors.New("failed to get result as vector")

	ErrResultEmpty = errors.New("result is empty")
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
