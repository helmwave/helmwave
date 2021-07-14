package release

import "helm.sh/helm/v3/pkg/action"

func (rel *Config) Rollback() error {
	client := action.NewRollback(rel.Cfg())
	return client.Run(rel.Name)

}
