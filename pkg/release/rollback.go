package release

import (
	"context"
	"fmt"

	"github.com/helmwave/helmwave/pkg/helper"
)

func (rel *config) Rollback(ctx context.Context, version int) (err error) {
	ctx = helper.ContextWithReleaseUniq(ctx, rel.Uniq())

	// Run hooks
	err = rel.Lifecycle.RunPreRollback(ctx)
	if err != nil {
		return
	}

	defer func() {
		lifecycleErr := rel.Lifecycle.RunPostRollback(ctx)
		if lifecycleErr != nil {
			rel.Logger().Errorf("got an error from postrollback hooks: %v", lifecycleErr)
			if err == nil {
				err = lifecycleErr
			}
		}
	}()

	client := rel.newRollback()

	if version > 0 {
		client.Version = version
		rel.Logger().Infof("Be careful! Rollback to %d revision", version)
	}

	err = client.Run(rel.Name())
	if err != nil {
		err = fmt.Errorf("failed to rollback release %s: %w", rel.Uniq(), err)
	}

	return
}
