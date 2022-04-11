package release

import (
	"fmt"

	"helm.sh/helm/v3/pkg/action"
)

func (rel *config) Rollback(version int) error {
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
