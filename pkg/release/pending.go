package release

import (
	"context"

	"github.com/invopop/jsonschema"
)

// PendingStrategy is a type for enumerating strategies for handling pending releases.
type PendingStrategy string

const (
	// PendingStrategyRollback rolls back pending release.
	PendingStrategyRollback PendingStrategy = "rollback"
	// PendingStrategyUninstall uninstalls pending release.
	PendingStrategyUninstall PendingStrategy = "uninstall"
)

func (PendingStrategy) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			PendingStrategyRollback,
			PendingStrategyUninstall,
			"",
		},
	}
}

func (rel *config) isPending() (bool, error) {
	status, err := rel.Status()
	if err != nil {
		return false, err
	}

	return status.Info.Status.IsPending(), nil
}

func (rel *config) fixPending(ctx context.Context) error {
	switch rel.PendingReleaseStrategy {
	case PendingStrategyRollback:
		return rel.Rollback(ctx, 0)
	case PendingStrategyUninstall:
		_, err := rel.Uninstall(ctx)

		return err
	default:
		return ErrPendingRelease
	}
}
