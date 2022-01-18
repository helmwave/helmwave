package release

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

func (rel *config) GetChart() (*chart.Chart, error) {
	// Hmm nice action bro
	client := rel.newInstall()

	ch, err := client.ChartPathOptions.LocateChart(rel.Chart().Name, rel.Helm())
	if err != nil {
		return nil, fmt.Errorf("failed to locate chart %s: %w", rel.Chart().Name, err)
	}

	c, err := loader.Load(ch)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart %s: %w", rel.Chart().Name, err)
	}

	if err := rel.chartCheck(c); err != nil {
		return nil, err
	}

	return c, nil
}

func (rel *config) chartCheck(ch *chart.Chart) error {
	if req := ch.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(ch, req); err != nil {
			return fmt.Errorf("failed to check chart %s dependencies: %w", ch.Name(), err)
		}
	}

	if !(ch.Metadata.Type == "" || ch.Metadata.Type == "application") {
		rel.Logger().Warnf("%s charts are not installable", ch.Metadata.Type)
	}

	if ch.Metadata.Deprecated {
		return fmt.Errorf("⚠️ Chart %s is deprecated", ch.Name())
	}

	return nil
}

func (rel *config) ChartDepsUpd() error {
	return rel.chartDepsUpd(rel.Chart().Name, rel.Helm())
}

func (rel *config) chartDepsUpd(name string, settings *helm.EnvSettings) error {
	if !helper.IsExists(filepath.Clean(name)) {
		rel.Logger().Info("skipping updating dependencies for remote chart")

		return nil
	}

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

	if err := man.Update(); err != nil {
		return fmt.Errorf("failed to update %s chart dependencies: %w", name, err)
	}

	return nil
}
