package hooks

import (
	"errors"
	"fmt"
)

var ErrUnknownFormat = errors.New("unknown format")

type CreatePipeError struct {
	Err error
}

func NewCreatePipeError(err error) error {
	return &CreatePipeError{Err: err}
}

func (err CreatePipeError) Error() string {
	return fmt.Sprintf("failed to create command pipe: %s", err.Err)
}

func (err CreatePipeError) Unwrap() error {
	return err.Err
}

type CommandRunError struct {
	Err error
}

func NewCommandRunError(err error) error {
	return &CommandRunError{Err: err}
}

func (err CommandRunError) Error() string {
	return fmt.Sprintf("failed to run command: %s", err.Err)
}

func (err CommandRunError) Unwrap() error {
	return err.Err
}

type CommandReadOutputError struct {
	Err error
}

func NewCommandReadOutputError(err error) error {
	return &CommandReadOutputError{Err: err}
}

func (err CommandReadOutputError) Error() string {
	return fmt.Sprintf("failed to read command stdout: %s", err.Err)
}

func (err CommandReadOutputError) Unwrap() error {
	return err.Err
}

type YAMLDecodeError struct {
	Err error
}

func NewYAMLDecodeError(err error) error {
	return &YAMLDecodeError{Err: err}
}

func (err YAMLDecodeError) Error() string {
	return fmt.Sprintf("failed to decode lifecycle config from YAML: %s", err.Err)
}

func (err YAMLDecodeError) Unwrap() error {
	return err.Err
}
