package release

import (
	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"os"
	"path/filepath"
)

func (rel *Config) upgrade (cfg *action.Configuration, helm *helm.EnvSettings) (*release.Release, error){
	client := action.NewUpgrade(cfg)
	// Merge
	err := mergo.Merge(client, rel)
	if err != nil {
		return nil, err
	}

	locateChart, err := client.ChartPathOptions.LocateChart(rel.Chart, helm)
	if err != nil {
		return nil, err
	}


	valOpts := &values.Options{ValueFiles: rel.Values}
	vals, err := valOpts.MergeValues(getter.All(helm))
	if err != nil {
		return nil, err
	}


	ch, err := loader.Load(locateChart)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if req := ch.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(ch, req); err != nil {
			return nil, err
		}
	}

	if !(ch.Metadata.Type == "" || ch.Metadata.Type == "application") {
		log.Warnf("%s charts are not installable \n", ch.Metadata.Type)
	}

	if ch.Metadata.Deprecated {
		log.Warn("⚠️ This locateChart is deprecated")
	}

	if !rel.isInstalled(cfg) {
		log.Debugf("🧐 Release %q in %q does not exist. Installing it now.", rel.ReleaseName, rel.Namespace)
		return rel.install(cfg, ch, vals)
	}

	return client.Run(rel.ReleaseName, ch, vals)

}

func (rel *Config) chartDepsUpd(settings *helm.EnvSettings) error {
	client := action.NewDependency()
	man := &downloader.Manager{
		Out:              os.Stdout,
		ChartPath:        filepath.Clean(rel.Chart),
		Keyring:          client.Keyring,
		SkipUpdate:       client.SkipRefresh,
		Getters:          getter.All(settings),
		RepositoryConfig: settings.RepositoryConfig,
		RepositoryCache:  settings.RepositoryCache,
		Debug:            settings.Debug,
	}
	if client.Verify {
		man.Verify = downloader.VerifyAlways
	}
	return man.Update()
}

