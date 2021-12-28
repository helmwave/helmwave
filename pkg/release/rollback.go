package release

import (
	"fmt"

	"helm.sh/helm/v3/pkg/action"
)

func (rel *config) Rollback() error {
	client := action.NewRollback(rel.Cfg())

	if err := client.Run(rel.Name()); err != nil {
		return fmt.Errorf("failed to rollback release %s: %w", rel.Uniq(), err)
	}

	return nil
}
