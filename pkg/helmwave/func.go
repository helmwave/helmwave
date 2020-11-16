package helmwave

import (
	"fmt"
	"github.com/zhilyaev/helmwave/pkg/yml"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"log"
	"os"
)

func (c *Config) ReadHelmWaveYml() {
	yml.Read(c.Yml.File, &c.Yml.Body)
	if c.Yml.Body.Version != c.Version {
		fmt.Println("⚠️ Unsupported version", c.Yml.Body.Version)
	}
}

func (c *Config) Log(format string, v ...interface{}) {
	if c.Debug {
		format = fmt.Sprintf("[debug] %s\n", format)
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

func (c Config) ActionCfg(ns string, settings *helm.EnvSettings) (*action.Configuration, error) {
	cfg := new(action.Configuration)
	helmDriver := os.Getenv("HELM_DRIVER")
	if len(ns) == 0 {
		ns = settings.Namespace()
	}

	err := cfg.Init(settings.RESTClientGetter(), ns, helmDriver, c.Log)
	return cfg, err
}
