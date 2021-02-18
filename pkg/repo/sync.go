package repo

import (
	log "github.com/sirupsen/logrus"
	helm "helm.sh/helm/v3/pkg/cli"
)

func Sync(repos []*Config, settings *helm.EnvSettings) (err error) {
	log.Info("🗄 Sync repositories")
	for _, r := range repos {
		err := r.Install(settings)
		if err != nil {
			return err
		}
	}

	return nil
}
