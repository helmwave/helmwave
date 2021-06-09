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
	"helm.sh/helm/v3/pkg/storage/driver"
	"os"
	"path/filepath"
)

func (rel *Config) Install(cfg *action.Configuration, settings *helm.EnvSettings) (*release.Release, error) {
	err := rel.chartDepsUpd(settings)
	if err != nil {
		log.Debug(err)
	}
	// I hate private field
	client := action.NewUpgrade(cfg)
	err = mergo.Merge(client, rel.Options)
	if err != nil {
		return nil, err
	}

	chart, err := client.ChartPathOptions.LocateChart(rel.Chart, settings)
	if err != nil {
		return nil, err
	}

	v := make([]string, len(rel.Values))
	for i := range rel.Values {
		v[i] = rel.Values[i].GetPath()
	}
	valOpts := &values.Options{ValueFiles: v}
	vals, err := valOpts.MergeValues(getter.All(settings))
	if err != nil {
		return nil, err
	}

	ch, err := loader.Load(chart)
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
		log.Warn("⚠️ This chart is deprecated")
	}

	if client.Install {
		// If a release does not exist, install it.
		histClient := action.NewHistory(cfg)
		histClient.Max = 1
		_, err := histClient.Run(rel.Name)
		if err == driver.ErrReleaseNotFound {
			log.Debugf("🧐 Release %q in %q does not exist. Installing it now.", rel.Name, rel.Options.Namespace)

			instClient := action.NewInstall(cfg)

			instClient.DryRun = client.DryRun

			instClient.CreateNamespace = true
			instClient.ReleaseName = rel.Name
			instClient.Namespace = client.Namespace

			// Mmm... Nice.
			instClient.ChartPathOptions = client.ChartPathOptions
			instClient.DisableHooks = client.DisableHooks
			instClient.SkipCRDs = client.SkipCRDs
			instClient.Timeout = client.Timeout
			instClient.Wait = client.Wait
			instClient.Devel = client.Devel
			instClient.Atomic = client.Atomic
			instClient.PostRenderer = client.PostRenderer
			instClient.DisableOpenAPIValidation = client.DisableOpenAPIValidation
			instClient.SubNotes = client.SubNotes
			instClient.Description = client.Description

			if instClient.DryRun {
				//instClient.ReleaseName = "RELEASE-NAME"
				instClient.Replace = true
			}

			return instClient.Run(ch, vals)

		} else if err != nil {
			return nil, err
		}

	}

	return client.Run(rel.Name, ch, vals)

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
