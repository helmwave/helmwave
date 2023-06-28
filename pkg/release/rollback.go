package release

import (
	"context"
	"fmt"

	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v3/pkg/action"
)

func (rel *config) Rollback(ctx context.Context, version int) error {
	ctx = helper.ContextWithReleaseUniq(ctx, rel.Uniq())

	// Run hooks
	err := rel.Lifecycle.RunPreRollback(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err := rel.Lifecycle.RunPostRollback(ctx)
		if err != nil {
			rel.Logger().Errorf("got an error from postrollback hooks: %v", err)
		}
	}()

	client := action.NewRollback(rel.Cfg())

	client.CleanupOnFail = rel.CleanupOnFail
	client.MaxHistory = rel.MaxHistory
	client.Recreate = rel.Recreate
	client.Timeout = rel.Timeout

	client.DisableHooks = rel.DisableHooks
	client.Timeout = rel.Timeout
	client.Wait = rel.Wait
	client.WaitForJobs = rel.WaitForJobs
	client.Force = rel.Force

	if version > 0 {
		client.Version = version
		rel.Logger().Infof("Be careful! Rollback to %d revision", version)
	}

	if err := client.Run(rel.Name()); err != nil {
		return fmt.Errorf("failed to rollback release %s: %w", rel.Uniq(), err)
	}

	return nil
}
