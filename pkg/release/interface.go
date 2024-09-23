package release

import (
	"context"
	"fmt"
	"slices"

	"github.com/helmwave/helmwave/pkg/fileref"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/hooks"
	"github.com/helmwave/helmwave/pkg/log"
	"github.com/helmwave/helmwave/pkg/monitor"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/invopop/jsonschema"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
)

// Config is an interface to manage particular helm release.
type Config interface {
	log.LoggerGetter
	HelmActionRunner

	Uniq() uniqname.UniqName
	AllowFailure() bool
	DryRun(dryRun bool)
	HideSecret(hideSecret bool)
	ChartDepsUpd() error
	DownloadChart(tmpDir string) error
	BuildValues(ctx context.Context, dir, templater string) error
	Name() string
	Namespace() string
	Chart() *Chart
	SetChartName(string)
	DependsOn() []*DependsOnReference
	SetDependsOn(deps []*DependsOnReference)
	Tags() []string
	Repo() string
	Values() []fileref.Config
	HelmWait() bool
	KubeContext() string
	Cfg() *action.Configuration
	HooksDisabled() bool
	OfflineKubeVersion() *chartutil.KubeVersion
	Validate() error
	Monitors() []MonitorReference
	NotifyMonitorsFailed(ctx context.Context, monitors ...monitor.Config)
	Lifecycle() hooks.Lifecycle
}

type HelmActionRunner interface {
	SyncDryRun(ctx context.Context, runHooks bool) (*release.Release, error)
	Sync(ctx context.Context, runHooks bool) (*release.Release, error)
	Uninstall(ctx context.Context) (*release.UninstallReleaseResponse, error)
	Get(version int) (*release.Release, error)
	List() (*release.Release, error)
	Rollback(ctx context.Context, version int) error
	Status() (*release.Release, error)
}

// Configs type of array Config.
type Configs []Config

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse YAML into `Config` interface.
func (r *Configs) UnmarshalYAML(node *yaml.Node) error {
	rr := make([]*config, 0)
	if err := node.Decode(&rr); err != nil {
		return fmt.Errorf("failed to decode release config from YAML: %w", err)
	}

	*r = helper.SlicesMap(rr, func(r *config) Config {
		r.buildAfterUnmarshal(rr)

		return r
	})

	return nil
}

func (Configs) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
	}
	var l []*config

	return r.Reflect(&l)
}

func (r Configs) Contains(rel Config) (Config, bool) {
	return r.ContainsUniq(rel.Uniq())
}

func (r Configs) ContainsUniq(uniq uniqname.UniqName) (Config, bool) {
	i := slices.IndexFunc(r, func(rel Config) bool {
		return rel.Uniq().Equal(uniq)
	})

	if i == -1 {
		return nil, false
	}

	return r[i], true
}
