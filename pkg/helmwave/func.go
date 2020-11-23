package helmwave

import (
	log "github.com/sirupsen/logrus"
	"github.com/zhilyaev/helmwave/pkg/yml"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"os"
)

func (c *Config) ReadHelmWaveYml() {
	yml.Read(c.Yml.File, &c.Yml.Body)
	if c.Yml.Body.Version != c.Version {
		log.Warn("⚠️ Unsupported version ", c.Yml.Body.Version)
		log.Debug("🌊 HelmWave version ", c.Version)
	}
}

func (c Config) ActionCfg(ns string, settings *helm.EnvSettings) (*action.Configuration, error) {
	cfg := new(action.Configuration)
	helmDriver := os.Getenv("HELM_DRIVER")
	if len(ns) == 0 {
		ns = settings.Namespace()
	}

	err := cfg.Init(settings.RESTClientGetter(), ns, helmDriver, log.Debugf)
	return cfg, err
}
