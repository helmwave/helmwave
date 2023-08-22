package log

import "fmt"

type InvalidLogLevelError struct {
	Err   error
	Level string
}

func NewInvalidLogLevelError(level string, err error) error {
	return &InvalidLogLevelError{Level: level, Err: err}
}

func (err InvalidLogLevelError) Error() string {
	return fmt.Sprintf("failed to parse log level %q: %s", err.Level, err.Err)
}

//nolint:errorlint
func (InvalidLogLevelError) Is(target error) bool {
	switch target.(type) {
	case InvalidLogLevelError, *InvalidLogLevelError:
		return true
	default:
		return false
	}
}

func (err InvalidLogLevelError) Unwrap() error {
	return err.Err
}
