package helper

import (
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"os"
)

func ActionCfg(ns string, settings *helm.EnvSettings) (*action.Configuration, error) {
	cfg := new(action.Configuration)
	helmDriver := os.Getenv("HELM_DRIVER")
	if len(ns) == 0 {
		ns = settings.Namespace()
	}

	err := cfg.Init(settings.RESTClientGetter(), ns, helmDriver, log.Debugf)
	return cfg, err
}
