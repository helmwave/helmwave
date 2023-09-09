package monitor

import (
	"errors"
	"fmt"
)

var ErrFailureStreak = errors.New("monitor triggered failure threshold")

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
