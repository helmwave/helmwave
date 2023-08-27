package release

import (
	"errors"
	"fmt"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"helm.sh/helm/v3/pkg/storage/driver"
)

var (
	ErrNameEmpty = errors.New("release name is empty")

	// ErrMissingDependency is returned when values can't be used and are skipped.
	ErrMissingDependency = errors.New("dependency is missing")

	// ErrPendingRelease is an error for fail strategy that release is in pending status.
	ErrPendingRelease = errors.New("release is in pending status")

	// ErrValuesNotExist is returned when values can't be used and are skipped.
	ErrValuesNotExist = errors.New("values file doesn't exist")

	// ErrNotFound is an error for not found release.
	ErrNotFound = driver.ErrReleaseNotFound

	// ErrFoundMultiple is an error for multiple releases found by name.
	ErrFoundMultiple = errors.New("found multiple releases o_0")

	// ErrDepFailed is an error thrown when dependency release fails.
	ErrDepFailed = errors.New("dependency failed")

	ErrUnknownFormat = errors.New("unknown format")
)

type DuplicateError struct {
	Uniq uniqname.UniqName
}

func NewDuplicateError(uniq uniqname.UniqName) error {
	return &DuplicateError{Uniq: uniq}
}

func (err DuplicateError) Error() string {
	return fmt.Sprintf("release duplicate: %s", err.Uniq.String())
}

func (DuplicateError) Is(target error) bool {
	switch target.(type) {
	case DuplicateError, *DuplicateError:
		return true
	default:
		return false
	}
}

type InvalidNamespaceError struct {
	Namespace string
}

func NewInvalidNamespaceError(namespace string) error {
	return &InvalidNamespaceError{Namespace: namespace}
}

func (err InvalidNamespaceError) Error() string {
	return fmt.Sprintf("invalid namespace: %s", err.Namespace)
}

func (InvalidNamespaceError) Is(target error) bool {
	switch target.(type) {
	case InvalidNamespaceError, *InvalidNamespaceError:
		return true
	default:
		return false
	}
}

type YAMLDecodeDependsOnError struct {
	Err       error
	DependsOn string
}

func NewYAMLDecodeDependsOnError(depends_on string, err error) error {
	return &YAMLDecodeDependsOnError{DependsOn: depends_on, Err: err}
}

func (err YAMLDecodeDependsOnError) Error() string {
	return fmt.Sprintf("failed to decode depends_on reference %q from YAML: %s", err.DependsOn, err.Err)
}

func (err YAMLDecodeDependsOnError) Unwrap() error {
	return err.Err
}

func (YAMLDecodeDependsOnError) Is(target error) bool {
	switch target.(type) {
	case YAMLDecodeDependsOnError, *YAMLDecodeDependsOnError:
		return true
	default:
		return false
	}
}
