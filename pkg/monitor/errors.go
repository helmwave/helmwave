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

type MonitorInitError struct {
	Err error
}

func NewMonitorInitError(err error) error {
	return &MonitorInitError{Err: err}
}

func (err MonitorInitError) Error() string {
	return fmt.Sprintf("monitor failed to initialize: %s", err.Err)
}

func (err MonitorInitError) Unwrap() error {
	return err.Err
}

func (MonitorInitError) Is(target error) bool {
	switch target.(type) {
	case MonitorInitError, *MonitorInitError:
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
	return fmt.Sprintf("failed to decode lifecycle config from YAML: %s", err.Err)
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

type DuplicateError struct {
	Name string
}

func NewDuplicateError(name string) error {
	return &DuplicateError{Name: name}
}

func (err DuplicateError) Error() string {
	return fmt.Sprintf("monitor duplicate: %s", err.Name)
}

func (DuplicateError) Is(target error) bool {
	switch target.(type) {
	case DuplicateError, *DuplicateError:
		return true
	default:
		return false
	}
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

func (NotExistsError) Is(target error) bool {
	switch target.(type) {
	case NotExistsError, *NotExistsError:
		return true
	default:
		return false
	}
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

func (SubMonitorError) Is(target error) bool {
	switch target.(type) {
	case SubMonitorError, *SubMonitorError:
		return true
	default:
		return false
	}
}
