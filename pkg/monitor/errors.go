package monitor

import (
	"errors"
	"fmt"
)

var (
	ErrNameEmpty = errors.New("release name is empty")

	ErrFailureStreak = errors.New("monitor triggered failure threshold")

	ErrLowTotalTimeout = errors.New("total timeout is less than iteration timeout")

	ErrLowInterval = errors.New("interval cannot be zero")
)

type InitError struct {
	Err error
}

func NewMonitorInitError(err error) error {
	return &InitError{Err: err}
}

func (err InitError) Error() string {
	return fmt.Sprintf("monitor failed to initialize: %s", err.Err)
}

func (err InitError) Unwrap() error {
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

type DuplicateError struct {
	Name string
}

func NewDuplicateError(name string) error {
	return &DuplicateError{Name: name}
}

func (err DuplicateError) Error() string {
	return fmt.Sprint("monitor duplicate: ", err.Name)
}

type NotExistsError struct {
	Name string
}

func NewNotExistsError(name string) error {
	return &NotExistsError{Name: name}
}

func (err NotExistsError) Error() string {
	return fmt.Sprintf("monitor doesn't exist: %s", err.Name)
}

type SubMonitorError struct {
	Err error
}

func NewSubMonitorError(err error) error {
	return &SubMonitorError{Err: err}
}

func (err SubMonitorError) Error() string {
	return fmt.Sprintf("submonitor config error: %s", err.Err)
}

func (err SubMonitorError) Unwrap() error {
	return err.Err
}
