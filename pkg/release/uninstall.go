package release

import (
	"context"
	"fmt"

	"github.com/helmwave/helmwave/pkg/helper"
	releaseiface "helm.sh/helm/v4/pkg/release"
)

func (rel *config) Uninstall(ctx context.Context) (resp *releaseiface.UninstallReleaseResponse, err error) {
	ctx = helper.ContextWithReleaseUniq(ctx, rel.Uniq())

	// Run hooks
	err = rel.LifecycleF.RunPreDown(ctx)
	if err != nil {
		return
	}

	defer func() {
		lifecycleErr := rel.LifecycleF.RunPostDown(ctx)
		if lifecycleErr != nil {
			if err == nil {
				err = lifecycleErr
			}
		}
	}()

	client := rel.newUninstall()

	resp, err = client.Run(rel.Name())
	if err != nil {
		err = fmt.Errorf("failed to uninstall release %s: %w", rel.Uniq(), err)
	}

	return
}
