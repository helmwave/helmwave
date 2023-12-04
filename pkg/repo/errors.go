package repo

import (
	"errors"
	"fmt"
)

var (
	ErrNameEmpty = errors.New("repository name is empty")
	ErrURLEmpty  = errors.New("repository url is empty")
)

type InvalidURLError struct {
	URL string
}

func NewInvalidURLError(url string) error {
	return &InvalidURLError{URL: url}
}

func (err InvalidURLError) Error() string {
	return fmt.Sprintf("invalid URL: %s", err.URL)
}

type DuplicateError struct {
	Name string
}

func NewDuplicateError(name string) error {
	return &DuplicateError{Name: name}
}

func (err DuplicateError) Error() string {
	return fmt.Sprintf("repository duplicate: %s", err.Name)
}

type NotFoundError struct {
	Name string
}

func NewNotFoundError(name string) error {
	return &NotFoundError{Name: name}
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("🗄 repository not found: %s", err.Name)
}
