package release

import (
	"context"
	"fmt"

	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *config) Uninstall(ctx context.Context) (*release.UninstallReleaseResponse, error) {
	ctx = helper.ContextWithReleaseUniq(ctx, rel.Uniq())

	// Run hooks
	err := rel.Lifecycle.RunPreDown(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rel.Lifecycle.RunPostDown(ctx)
		if err != nil {
			rel.Logger().Errorf("got an error from postdown hooks: %v", err)
		}
	}()

	client := rel.newUninstall()

	resp, err := client.Run(rel.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to uninstall release %s: %w", rel.Uniq(), err)
	}

	return resp, nil
}
