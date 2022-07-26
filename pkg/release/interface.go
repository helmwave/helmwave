package release

import (
	"fmt"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/release"
)

// Config is an interface to manage particular helm release.
type Config interface {
	Uniq() uniqname.UniqName
	Sync() (*release.Release, error)
	AllowFailure() bool
	DryRun(bool)
	ChartDepsUpd() error
	In([]Config) bool
	BuildValues(string, string) error
	Uninstall() (*release.UninstallReleaseResponse, error)
	Get() (*release.Release, error)
	List() (*release.Release, error)
	Rollback(int) error
	Status() (*release.Release, error)
	Name() string
	Namespace() string
	Chart() Chart
	DependsOn() []string
	Tags() []string
	Repo() string
	Values() []ValuesReference
	Logger() *log.Entry
}

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse YAML into `Config` interface.
func UnmarshalYAML(node *yaml.Node) ([]Config, error) {
	r := make([]*config, 0)
	if err := node.Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode release config from YAML: %w", err)
	}

	res := make([]Config, len(r))
	for i := range r {
		res[i] = r[i]
	}

	return res, nil
}
