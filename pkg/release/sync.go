package release

import (
	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *config) Sync() (*release.Release, error) {
	// DependsON
	if err := rel.waitForDependencies(); err != nil {
		return nil, err
	}

	return rel.upgrade()
}

func (rel *config) Cfg() *action.Configuration {
	if rel.cfg == nil {
		var err error
		rel.cfg, err = helper.NewCfg(rel.Namespace())
		if err != nil {
			log.Fatal(err)

			return nil
		}
	}

	return rel.cfg
}

func (rel *config) Helm() *helm.EnvSettings {
	if rel.helm == nil {
		var err error
		rel.helm, err = helper.NewHelm(rel.Namespace())
		if err != nil {
			log.Fatal(err)

			return nil
		}

		rel.helm.Debug = helper.Helm.Debug
	}

	return rel.helm
}
