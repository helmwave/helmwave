package helmwave

import (
	log "github.com/sirupsen/logrus"
	"github.com/zhilyaev/helmwave/pkg/template"
	"github.com/zhilyaev/helmwave/pkg/yml"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"os"
)

func (c *Config) ReadBody(filepath string, destination *yml.Body) error {
	log.Info("Looking for ", filepath)
	err := yml.Read(filepath, destination)
	if err != nil {
		return err
	}

	return c.CheckVersion(c.Yml.Body.Version)
}

func (c *Config) CheckVersion(version string) error {
	if version != c.Version {
		log.Warn("⚠️ Unsupported version ", version)
		log.Debug("🌊 HelmWave version ", c.Version)
	}

	return nil
}

func (c *Config) ReadHelmWaveYml() error {
	return c.ReadBody(c.Yml.File, &c.Yml.Body)
}

func (c *Config) ReadHelmWavePlan() error {
	return c.ReadBody(c.Plan.File, &c.Plan.Body)
}

func (c *Config) RenderHelmWaveYml() error {
	return template.Tpl2yml(c.Tpl.File, c.Yml.File, nil)
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
