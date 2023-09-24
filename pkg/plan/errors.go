package plan

import (
	"errors"
	"fmt"
)

var (
	// ErrValidateFailed is returned for failed values validation.
	ErrValidateFailed = errors.New("validate failed")

	// ErrPlansAreTheSame is returned when trying to compare plan with itself.
	ErrPlansAreTheSame = errors.New("plan1 and plan2 are the same")

	ErrMultipleKubecontexts = errors.New("kubedog can't work with releases in multiple kubecontexts")

	// ErrManifestDirEmpty is an error for empty manifest dir.
	ErrManifestDirEmpty = fmt.Errorf("manifests dir %s is empty", Manifest)

	// ErrDeploy is returned when deploy is failed for whatever reason.
	ErrDeploy = errors.New("deploy failed")

	ErrInvalidPlandir = errors.New("filesystem not supported")
)

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
