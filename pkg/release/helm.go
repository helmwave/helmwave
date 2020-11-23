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
	"helm.sh/helm/v3/pkg/storage/driver"
	"os"
	"path/filepath"
)

func (rel *Config) DependencyUpdate(settings *helm.EnvSettings) error {
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

func (rel *Config) Sync(cfg *action.Configuration, settings *helm.EnvSettings) error {
	err := rel.DependencyUpdate(settings)
	if err != nil {
		log.Debug(err)
	}
	// I hate private field
	client := action.NewUpgrade(cfg)
	err = mergo.Merge(client, rel.Options)
	if err != nil {
		return err
	}

	chart, err := client.ChartPathOptions.LocateChart(rel.Chart, settings)
	if err != nil {
		return err
	}

	valOpts := &values.Options{ValueFiles: rel.Values}
	vals, err := valOpts.MergeValues(getter.All(settings))
	if err != nil {
		return err
	}

	ch, err := loader.Load(chart)
	if err != nil {
		log.Fatal(err)
	}

	if req := ch.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(ch, req); err != nil {
			return err
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
			log.Debugf("🧐 Release %q in %q does not exist. Installing it now.\n", rel.Name, rel.Options.Namespace)

			instClient := action.NewInstall(cfg)

			instClient.CreateNamespace = true
			instClient.ReleaseName = rel.Name
			instClient.Namespace = client.Namespace

			// Mmm... Nice.
			instClient.ChartPathOptions = client.ChartPathOptions
			instClient.DryRun = client.DryRun
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

			_, err := instClient.Run(ch, vals)
			return err

		} else if err != nil {
			return err
		}
	}

	_, err = client.Run(rel.Name, ch, vals)
	if err != nil {
		return err
	}
	return nil
}
