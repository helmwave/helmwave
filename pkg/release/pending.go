package release

import (
	"context"
	"errors"

	"helm.sh/helm/v3/pkg/release"
)

// PendingStrategy is a type for enumerating strategies for handling pending releases.
type PendingStrategy string

const (
	// PendingStrategyRollback rolls back pending release.
	PendingStrategyRollback PendingStrategy = "rollback"
	// PendingStrategyUninstall uninstalls pending release.
	PendingStrategyUninstall PendingStrategy = "uninstall"
)

// ErrPendingRelease is an error for fail strategy that release is in pending status.
var ErrPendingRelease = errors.New("release is in pending status")

func (rel *config) isPending() (bool, error) {
	status, err := rel.Status()
	if err != nil {
		return false, err
	}

	switch status.Info.Status {
	case release.StatusPendingInstall, release.StatusPendingRollback, release.StatusPendingUpgrade:
		return true, nil
	default:
		return false, nil
	}
}

func (rel *config) fixPending(ctx context.Context) error {
	switch rel.PendingReleaseStrategy {
	case PendingStrategyRollback:
		return rel.Rollback(0)
	case PendingStrategyUninstall:
		_, err := rel.Uninstall(ctx)

		return err
	default:
		return ErrPendingRelease
	}
}
