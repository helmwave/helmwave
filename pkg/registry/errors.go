package registry

import (
	"errors"
	"fmt"
)

var ErrNameEmpty = errors.New("registry name is empty")

type DuplicateError struct {
	Host string
}

func (err DuplicateError) Error() string {
	return fmt.Sprintf("registry duplicate: %s", err.Host)
}

//nolint:errorlint
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

func (err NotFoundError) Error() string {
	return fmt.Sprintf("🗄 registry not found: %s", err.Host)
}

//nolint:errorlint
func (NotFoundError) Is(target error) bool {
	switch target.(type) {
	case NotFoundError, *NotFoundError:
		return true
	default:
		return false
	}
}
