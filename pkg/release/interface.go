package release

import (
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/release"
)

// Config is an interface to manage particular helm release.
type Config interface {
	Uniq() uniqname.UniqName
	HandleDependencies([]Config)
	Sync() (*release.Release, error)
	NotifySuccess()
	NotifyFailed()
	DryRun(bool)
	ChartDepsUpd() error
	In([]Config) bool
	BuildValues(string, string) error
	Uninstall() (*release.UninstallReleaseResponse, error)
	Get() (*release.Release, error)
	List() (*release.Release, error)
	Rollback() error
	Status() (*release.Release, error)

	Name() string
	Namespace() string
	Chart() Chart
	DependsOn() []string
	Tags() []string
	Repo() string
	Values() []ValuesReference
}

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse YAML into `Config` interface.
func UnmarshalYAML(node *yaml.Node) ([]Config, error) {
	r := make([]*config, 0)
	if err := node.Decode(&r); err != nil {
		return nil, err
	}

	res := make([]Config, len(r))
	for i := range r {
		res[i] = r[i]
	}

	return res, nil
}
