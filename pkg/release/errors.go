package release

import (
	"errors"
	"fmt"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"helm.sh/helm/v3/pkg/storage/driver"
)

var (
	ErrNameEmpty = errors.New("release name is empty")

	// ErrPendingRelease is an error for fail strategy that release is in pending status.
	ErrPendingRelease = errors.New("release is in pending status")

	// ErrNotFound is an error for not found release.
	ErrNotFound = driver.ErrReleaseNotFound

	// ErrFoundMultiple is an error for multiple releases found by name.
	ErrFoundMultiple = errors.New("found multiple releases o_0")

	// ErrDepFailed is an error thrown when dependency release fails.
	ErrDepFailed = errors.New("dependency failed")

	ErrUnknownFormat = errors.New("unknown format")

	ErrDigestNotMatch = errors.New("chart digest doesn't match")
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

type InvalidNamespaceError struct {
	Namespace string
}

func NewInvalidNamespaceError(namespace string) error {
	return &InvalidNamespaceError{Namespace: namespace}
}

func (err InvalidNamespaceError) Error() string {
	return fmt.Sprintf("invalid namespace: %s", err.Namespace)
}

type YAMLDecodeDependsOnError struct {
	Err       error
	DependsOn string
}

func NewYAMLDecodeDependsOnError(dependsOn string, err error) error {
	return &YAMLDecodeDependsOnError{DependsOn: dependsOn, Err: err}
}

func (err YAMLDecodeDependsOnError) Error() string {
	return fmt.Sprintf("failed to decode depends_on reference %q from YAML: %s", err.DependsOn, err.Err)
}

func (err YAMLDecodeDependsOnError) Unwrap() error {
	return err.Err
}

type ChartCacheError struct {
	Err error
}

func NewChartCacheError(err error) error {
	return &ChartCacheError{Err: err}
}

func (err ChartCacheError) Error() string {
	return fmt.Sprintf("failed to find chart in helm cache: %s", err.Err)
}

func (err ChartCacheError) Unwrap() error {
	return err.Err
}

type HelmTestsError struct {
	Err error
}

func NewHelmTestsError(err error) error {
	return &HelmTestsError{Err: err}
}

func (err HelmTestsError) Error() string {
	return fmt.Sprintf("helm tests failed: %s", err.Err)
}

func (err HelmTestsError) Unwrap() error {
	return err.Err
}
