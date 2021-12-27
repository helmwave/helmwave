package helper

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var Helm = helm.New()

func NewCfg(ns string) (*action.Configuration, error) {
	cfg := new(action.Configuration)
	helmDriver := os.Getenv("HELM_DRIVER")
	config := genericclioptions.NewConfigFlags(false)
	config.Namespace = &ns
	err := cfg.Init(config, ns, helmDriver, log.Debugf)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func NewHelm(ns string) (*helm.EnvSettings, error) {
	env := helm.New()
	fs := &pflag.FlagSet{}
	env.AddFlags(fs)
	flag := fs.Lookup("namespace")

	if err := flag.Value.Set(ns); err != nil {
		return nil, err
	}

	return env, nil
}
