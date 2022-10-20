package release

import (
	"context"
	"fmt"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/log"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"gopkg.in/yaml.v3"
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
	BuildValues(string, string) error
	Uninstall(context.Context) (*release.UninstallReleaseResponse, error)
	Get() (*release.Release, error)
	List() (*release.Release, error)
	Rollback(int) error
	Status() (*release.Release, error)
	Name() string
	Namespace() string
	Chart() Chart
	DependsOn() []*DependsOnReference
	Tags() []string
	Repo() string
	Values() []ValuesReference
	HelmWait() bool
}

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse YAML into `Config` interface.
func UnmarshalYAML(node *yaml.Node) ([]Config, error) {
	r := make([]*config, 0)
	if err := node.Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode release config from YAML: %w", err)
	}

	res := make([]Config, len(r))
	for i := range r {
		r[i].buildAfterUnmarshal(r)
		res[i] = r[i]
	}

	return res, nil
}
