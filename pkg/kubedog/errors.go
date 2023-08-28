package kubedog

import "fmt"

type ParseError struct {
	err   error
	typ   string
	value string
}

func NewParseError(typ, value string, err error) error {
	return &ParseError{typ: typ, value: value, err: err}
}

func (err ParseError) Error() string {
	return fmt.Sprintf("failed to parse %q as %s: %s", err.value, err.typ, err.err)
}

func (err ParseError) Unwrap() error {
	return err.err
}

func (ParseError) Is(target error) bool {
	switch target.(type) {
	case ParseError, *ParseError:
		return true
	default:
		return false
	}
}

type InvalidValueError[T ~string] struct {
	value      T
	annotation string
	choices    []T
}

func NewInvalidValueError[T ~string](annotation string, value T, choices []T) error {
	return &InvalidValueError[T]{annotation: annotation, value: value, choices: choices}
}

func (err InvalidValueError[T]) Error() string {
	return fmt.Sprintf(
		"invalid value %q for annotation %q, should be one of %q",
		err.value,
		err.annotation,
		err.choices,
	)
}

func (InvalidValueError[T]) Is(target error) bool {
	switch target.(type) {
	case InvalidValueError[T], *InvalidValueError[T]:
		return true
	default:
		return false
	}
}

type EmptyContainerNameError struct {
	annotation string
	value      string
}

func NewEmptyContainerNameError(annotation, value string) error {
	return &EmptyContainerNameError{annotation: annotation, value: value}
}

func (err EmptyContainerNameError) Error() string {
	return fmt.Sprintf(
		"cannot parse %q in annotation %q as comma-separated list of non-empty container names",
		err.value,
		err.annotation,
	)
}

func (EmptyContainerNameError) Is(target error) bool {
	switch target.(type) {
	case EmptyContainerNameError, *EmptyContainerNameError:
		return true
	default:
		return false
	}
}
