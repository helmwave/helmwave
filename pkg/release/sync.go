package release

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/wayt/parallel"
	"github.com/zhilyaev/helmwave/pkg/helper"
	helm "helm.sh/helm/v3/pkg/cli"
	"os"
)

var emptyReleases = errors.New("releases are empty")

func (rel *Config) Sync(manifestPath string) error {
	log.Infof("🛥 %s", rel.UniqName())

	// I hate Private
	_ = os.Setenv("HELM_NAMESPACE", rel.Options.Namespace)
	settings := helm.New()
	cfg, err := helper.ActionCfg(rel.Options.Namespace, settings)
	if err != nil {
		return err
	}

	install, err := rel.Install(cfg, settings)
	if install != nil {
		log.Trace(install.Manifest)
	}

	if err != nil {
		return err
	}

	m := manifestPath + rel.UniqName() + ".yml"
	f, err := helper.CreateFile(m)
	if err != nil {
		return err
	}
	_, err = f.WriteString(install.Manifest)
	if err != nil {
		return err
	}

	return f.Close()
}

func (rel *Config) SyncWithFails(fails *[]*Config, manifestPath string) {
	err := rel.Sync(manifestPath)
	if err != nil {
		log.Error("❌ ", err)
		*fails = append(*fails, rel)
	} else {
		log.Infof("✅ %s", rel.UniqName())
	}
}

func Sync(releases []*Config, manifestPath string, async bool) (err error) {
	if len(releases) == 0 {
		return emptyReleases
	}

	log.Info("🛥 Sync releases")
	var fails []*Config

	if async {
		g := &parallel.Group{}
		log.Debug("🐞 Run in parallel mode")
		for i := range releases {
			g.Go(releases[i].SyncWithFails, &fails, manifestPath)
		}
		err := g.Wait()
		if err != nil {
			return err
		}
	} else {
		for _, r := range releases {
			r.SyncWithFails(&fails, manifestPath)
		}
	}

	n := len(releases)
	k := len(fails)

	log.Infof("Success %d / %d", n-k, n)
	if k > 0 {
		for _, rel := range fails {
			log.Errorf("%q was not deployed to %q", rel.Name, rel.Options.Namespace)
		}

		return errors.New("deploy failed")
	}
	return nil
}
