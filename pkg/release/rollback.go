package release

import "helm.sh/helm/v3/pkg/action"

func (rel *config) Rollback() error {
	client := action.NewRollback(rel.Cfg())

	return client.Run(rel.Name())
}
