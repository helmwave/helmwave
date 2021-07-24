package release

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

func (rel *Config) GetChart() (*chart.Chart, error) {
	// Hmm nice action bro
	client := rel.newInstall()

	ch, err := client.ChartPathOptions.LocateChart(rel.Chart.Name, rel.Helm())
	if err != nil {
		return nil, err
	}

	c, err := loader.Load(ch)
	if err != nil {
		return nil, err
	}

	if err := chartCheck(c); err != nil {
		return nil, err
	}

	return c, nil
}

func chartCheck(ch *chart.Chart) error {
	if req := ch.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(ch, req); err != nil {
			return err
		}
	}

	if !(ch.Metadata.Type == "" || ch.Metadata.Type == "application") {
		log.Warnf("%s charts are not installable \n", ch.Metadata.Type)
	}

	if ch.Metadata.Deprecated {
		return errors.New("⚠️ This locateChart is deprecated")
	}

	return nil
}

func (rel *Config) ChartDepsUpd() error {
	return chartDepsUpd(rel.Chart.Name, rel.Helm())
}

func chartDepsUpd(name string, settings *helm.EnvSettings) error {
	client := action.NewDependency()
	man := &downloader.Manager{
		Out:              io.Discard,
		ChartPath:        filepath.Clean(name),
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

	if settings.Debug {
		man.Out = os.Stdout
	}

	return man.Update()
}
