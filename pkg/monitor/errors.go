package monitor

import "fmt"

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
