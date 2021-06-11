package release

import (
	"errors"
	"os"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	log "github.com/sirupsen/logrus"
	helm "helm.sh/helm/v3/pkg/cli"
)

var emptyReleases = errors.New("releases are empty")

func Sync(releases []*Config, manifestPath string, flagParallel bool) (err error) {
	if len(releases) == 0 {
		return emptyReleases
	}

	log.Info("🛥 Sync releases")
	var fails []*Config

	if flagParallel {
		wg := parallel.NewWaitGroup()
		log.Debug("🐞 Run in parallel mode")
		wg.Add(len(releases))
		for i := range releases {
			go func(wg *parallel.WaitGroup, release *Config, fails *[]*Config, manifestPath string) {
				defer wg.Done()
				release.SyncWithFails(fails, manifestPath)
			}(wg, releases[i], &fails, manifestPath)
		}
		err := wg.Wait()
		if err != nil {
			return err
		}
	} else {
		for _, r := range releases {
			r.SyncWithFails(&fails, manifestPath)
		}
	}

	n := len(releases)

	return showSuccess(n, fails)
}

func showSuccess(n int, fails []*Config) error {
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

func (rel *Config) Sync(manifestPath string) error {
	log.Infof("🛥 %s", rel.UniqName())

	if err := rel.waitForDependencies(); err != nil {
		return err
	}

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
		rel.NotifyFailed()
		*fails = append(*fails, rel)
	} else {
		rel.NotifySuccess()
		log.Infof("✅ %s", rel.UniqName())
	}
}
