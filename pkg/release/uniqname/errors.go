package uniqname

import (
	"fmt"
)

type ValidationError struct {
	Uniq string
}

func (err ValidationError) Error() string {
	return fmt.Sprintf("failed to validate uniqname: %s", err.Uniq)
}

//nolint:errorlint
func (ValidationError) Is(target error) bool {
	switch target.(type) {
	case ValidationError, *ValidationError:
		return true
	default:
		return false
	}
}
