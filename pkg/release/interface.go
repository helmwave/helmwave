package release

import (
	"context"
	"fmt"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/log"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/invopop/jsonschema"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
)

// Config is an interface to manage particular helm release.
type Config interface {
	helper.EqualChecker[Config]
	log.LoggerGetter
	Uniq() uniqname.UniqName
	Sync(context.Context) (*release.Release, error)
	SyncDryRun(context.Context) (*release.Release, error)
	AllowFailure() bool
	DryRun(bool)
	ChartDepsUpd() error
	DownloadChart(string) error
	BuildValues(string, string) error
	Uninstall(context.Context) (*release.UninstallReleaseResponse, error)
	Get() (*release.Release, error)
	List() (*release.Release, error)
	Rollback(context.Context, int) error
	Status() (*release.Release, error)
	Name() string
	Namespace() string
	Chart() *Chart
	SetChartName(string)
	DependsOn() []*DependsOnReference
	Tags() []string
	Repo() string
	Values() []ValuesReference
	HelmWait() bool
	KubeContext() string
	Cfg() *action.Configuration
	HooksDisabled() bool
	OfflineKubeVersion() *chartutil.KubeVersion
}

// Configs type of array Config.
type Configs []Config

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse YAML into `Config` interface.
func (r *Configs) UnmarshalYAML(node *yaml.Node) error {
	rr := make([]*config, 0)
	if err := node.Decode(&rr); err != nil {
		return fmt.Errorf("failed to decode release config from YAML: %w", err)
	}

	*r = make([]Config, len(rr))
	for i := range rr {
		rr[i].buildAfterUnmarshal(rr)
		(*r)[i] = rr[i]
	}

	return nil
}

func (Configs) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{DoNotReference: true}
	var l []*config

	return r.Reflect(&l)
}
