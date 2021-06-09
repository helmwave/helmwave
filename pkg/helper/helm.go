package helper

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
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

func setNS(ns string) (*helm.EnvSettings, error) {
	env := helm.New()
	fs := &pflag.FlagSet{}
	env.AddFlags(fs)
	flag := fs.Lookup("namespace")
	err := flag.Value.Set(ns)
	if err != nil {
		return nil, err
	}

	return env, nil
}
