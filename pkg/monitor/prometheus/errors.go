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

type PrometheusClientError struct {
	Err error
}

func NewPrometheusClientError(err error) error {
	return &PrometheusClientError{Err: err}
}

func (err PrometheusClientError) Error() string {
	return fmt.Sprintf("failed to create prometheus client: %s", err.Err)
}

func (err PrometheusClientError) Unwrap() error {
	return err.Err
}

func (PrometheusClientError) Is(target error) bool {
	switch target.(type) {
	case PrometheusClientError, *PrometheusClientError:
		return true
	default:
		return false
	}
}
