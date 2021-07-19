package release

import (
	"os"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func (rel *Config) Sync() (*release.Release, error) {
	// DependsON
	if err := rel.waitForDependencies(); err != nil {
		return nil, err
	}

	return rel.upgrade()
}

func (rel *Config) newCfg() (*action.Configuration, error) {
	cfg := new(action.Configuration)
	helmDriver := os.Getenv("HELM_DRIVER")
	err := cfg.Init(genericclioptions.NewConfigFlags(false), rel.Namespace, helmDriver, log.Debugf)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (rel *Config) Cfg() *action.Configuration {
	if rel.cfg == nil {
		var err error
		rel.cfg, err = rel.newCfg()
		if err != nil {
			log.Fatal(err)
			return nil
		}
	}

	return rel.cfg
}

// func (rel *Config) helm() (*helm.EnvSettings, error) {
//	env := helm.New()
//	fs := &pflag.FlagSet{}
//	env.AddFlags(fs)
//	flag := fs.Lookup("namespace")
//	err := flag.Value.Set(rel.Namespace)
//	if err != nil {
//		return nil, err
//	}
//
//	return env, nil
// }
