package release

import (
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
)



func (rel *Config) Sync () (*release.Release,error) {
	cfg, err := rel.cfg()
	if err != nil {
		return nil, err
	}

	helmClient, err := rel.helm()
	if err != nil {
		return nil, err
	}

	return rel.upgrade(cfg, helmClient)
}

func (rel* Config) SyncAndSaveManifest() error {
	r, err := rel.Sync()
	if r != nil {
		log.Trace(r.Manifest)
	}

	if err != nil {
		return err
	}

	m := plan.PlanManifest + rel.UniqName() + ".yml"

	f, err := helper.CreateFile(m)
	if err != nil {
		return err
	}
	_, err = f.WriteString(r.Manifest)
	if err != nil {
		return err
	}

	return f.Close()
}



func (rel * Config) cfg() (*action.Configuration, error) {
	cfg := new(action.Configuration)
	helmDriver := os.Getenv("HELM_DRIVER")
	err := cfg.Init(genericclioptions.NewConfigFlags(false), rel.Namespace, helmDriver, log.Debugf)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (rel * Config) helm() (*helm.EnvSettings, error) {
	env := helm.New()
	fs := &pflag.FlagSet{}
	env.AddFlags(fs)
	flag := fs.Lookup("namespace")
	err := flag.Value.Set(rel.Namespace)
	if err != nil {
		return nil, err
	}

	return env, nil
}