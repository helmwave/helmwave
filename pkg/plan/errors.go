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
