package release

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

// Chart is structure for chart download options.
//
//nolint:lll
type Chart struct {
	action.ChartPathOptions `yaml:",inline" json:",inline"`
	Name                    string `yaml:"name" json:"name" jsonschema:"description=Name of the chart,example=bitnami/nginx,example=oci://ghcr.io/helmwave/unit-test-oci"`
}

// UnmarshalYAML flexible config.
func (u *Chart) UnmarshalYAML(node *yaml.Node) error {
	type raw Chart
	var err error

	switch node.Kind {
	case yaml.ScalarNode, yaml.AliasNode:
		err = node.Decode(&(u.Name))
	case yaml.MappingNode:
		err = node.Decode((*raw)(u))
	default:
		err = fmt.Errorf("unknown format")
	}

	if err != nil {
		return fmt.Errorf("failed to decode chart %q from YAML at %d line: %w", node.Value, node.Line, err)
	}

	return nil
}

func (u Chart) IsRemote() bool {
	return !helper.IsExists(filepath.Clean(u.Name))
}

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
		rrel.Logger().Warnf("⚠️ Chart %s is deprecated. Please update your chart.", ch.Name())
	}

	return nil
}

func (rel *config) ChartDepsUpd() error {
	if rel.Chart().IsRemote() {
		rel.Logger().Info("skipping updating dependencies for remote chart")

		return nil
	}

	return chartDepsUpd(rel.Chart().Name, rel.Helm())
}

func chartDepsUpd(name string, settings *helm.EnvSettings) error {
	client := action.NewDependency()
	man := &downloader.Manager{
		Out:              log.StandardLogger().Writer(),
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

	if err := man.Update(); err != nil {
		return fmt.Errorf("failed to update %s chart dependencies: %w", name, err)
	}

	return nil
}

func (rel *config) DownloadChart(tmpDir string) error {
	if !rel.Chart().IsRemote() {
		rel.Logger().Info("chart is local, skipping exporting it")

		return nil
	}

	pull := action.NewPullWithOpts(action.WithConfig(rel.Cfg()))
	pull.Settings = rel.Helm()
	rel.copyChartPathOptions(&pull.ChartPathOptions)

	pull.DestDir = path.Join(tmpDir, "charts", rel.Uniq().String())
	err := os.MkdirAll(pull.DestDir, 0o750)
	if err != nil {
		return fmt.Errorf("failed to create temporary directory for chart: %w", err)
	}

	logs, err := pull.Run(rel.Chart().Name)
	if logs != "" {
		log.StandardLogger().Print(logs)
	}

	if err != nil {
		return fmt.Errorf("failed to download and unarchive chart: %w", err)
	}

	return nil
}

func (rel *config) SetChart(name string) {
	rel.lock.Lock()
	rel.ChartF.Name = name
	rel.lock.Unlock()
}
