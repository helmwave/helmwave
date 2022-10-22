package release

import (
	"context"
	"fmt"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/log"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/invopop/jsonschema"
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
	DependsOn() []uniqname.UniqName
	Tags() []string
	Repo() string
	Values() []ValuesReference
	HelmWait() bool
}

// Configs type of array Config.
type Configs []Config

// UnmarshalYAML is an unmarshaller for github.com/goccy/go-yaml to parse YAML into `Config` interface.
func (r *Configs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	rr := make([]*config, 0)
	if err := unmarshal(&rr); err != nil {
		return fmt.Errorf("failed to decode registry config from YAML: %w", err)
	}

	*r = make([]Config, len(rr))
	for i := range rr {
		rr[i].buildAfterUnmarshal()
		(*r)[i] = rr[i]
	}

	return nil
}

func (Configs) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{DoNotReference: true}
	var l []*config

	return r.Reflect(&l)
}
