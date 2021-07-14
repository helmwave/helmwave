package release

import "helm.sh/helm/v3/pkg/action"

func (rel *Config) Rollback() error {
	var err error
	rel.cfg, err = rel.newCfg()
	if err != nil {
		return err
	}

	client := action.NewRollback(rel.cfg)
	return client.Run(rel.Name)

}
