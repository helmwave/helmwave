package registry

import (
	"errors"
	"fmt"
)

var ErrNameEmpty = errors.New("registry name is empty")

type DuplicateError struct {
	Host string
}

func NewDuplicateError(host string) error {
	return &DuplicateError{Host: host}
}

func (err DuplicateError) Error() string {
	return fmt.Sprintf("registry duplicate: %s", err.Host)
}

func (DuplicateError) Is(target error) bool {
	switch target.(type) {
	case DuplicateError, *DuplicateError:
		return true
	default:
		return false
	}
}

type NotFoundError struct {
	Host string
}

func NewNotFoundError(host string) error {
	return &NotFoundError{Host: host}
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("🗄 registry not found: %s", err.Host)
}

func (NotFoundError) Is(target error) bool {
	switch target.(type) {
	case NotFoundError, *NotFoundError:
		return true
	default:
		return false
	}
}

type LoginError struct {
	Err error
}

func NewLoginError(err error) error {
	return &LoginError{Err: err}
}

func (err LoginError) Error() string {
	return fmt.Sprintf("failed to login in helm registry: %s", err.Err)
}

func (err LoginError) Unwrap() error {
	return err.Err
}

func (LoginError) Is(target error) bool {
	switch target.(type) {
	case LoginError, *LoginError:
		return true
	default:
		return false
	}
}

type YAMLDecodeError struct {
	Err error
}

func NewYAMLDecodeError(err error) error {
	return &YAMLDecodeError{Err: err}
}

func (err YAMLDecodeError) Error() string {
	return fmt.Sprintf("failed to decode registry config from YAML: %s", err.Err)
}

func (err YAMLDecodeError) Unwrap() error {
	return err.Err
}

func (YAMLDecodeError) Is(target error) bool {
	switch target.(type) {
	case YAMLDecodeError, *YAMLDecodeError:
		return true
	default:
		return false
	}
}
