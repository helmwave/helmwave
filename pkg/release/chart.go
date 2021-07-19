package release

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

func (rel *Config) GetChart() (*chart.Chart, error) {
	// Hmm nice action bro
	client := rel.newInstall()

	ch, err := client.ChartPathOptions.LocateChart(rel.Chart.Name, helper.Helm)
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
	return chartDepsUpd(rel.Chart.Name)
}

func chartDepsUpd(name string) error {
	client := action.NewDependency()
	man := &downloader.Manager{
		Out:              io.Discard,
		ChartPath:        filepath.Clean(name),
		Keyring:          client.Keyring,
		SkipUpdate:       client.SkipRefresh,
		Getters:          getter.All(helper.Helm),
		RepositoryConfig: helper.Helm.RepositoryConfig,
		RepositoryCache:  helper.Helm.RepositoryCache,
		Debug:            helper.Helm.Debug,
	}
	if client.Verify {
		man.Verify = downloader.VerifyAlways
	}

	if helper.Helm.Debug {
		man.Out = os.Stdout
	}

	return man.Update()
}
