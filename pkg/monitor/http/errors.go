package http

import (
	"errors"
	"fmt"
)

var ErrURLEmpty = errors.New("URL cannot be empty")

type RequestError struct {
	Err error
}

func NewRequestError(err error) error {
	return &RequestError{Err: err}
}

func (err RequestError) Error() string {
	return fmt.Sprintf("failed to create HTTP request: %s", err.Err)
}

func (err RequestError) Unwrap() error {
	return err.Err
}

type ResponseError struct {
	Err error
}

func NewResponseError(err error) error {
	return &ResponseError{Err: err}
}

func (err ResponseError) Error() string {
	return fmt.Sprintf("failed to get HTTP response: %s", err.Err)
}

func (err ResponseError) Unwrap() error {
	return err.Err
}

type UnexpectedStatusError struct {
	StatusCode int
}

func NewUnexpectedStatusError(status int) error {
	return &UnexpectedStatusError{StatusCode: status}
}

func (err UnexpectedStatusError) Error() string {
	return fmt.Sprintf("unexpected status code %d", err.StatusCode)
}
