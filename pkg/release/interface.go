package release

import (
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/template"
	"helm.sh/helm/v3/pkg/release"
)

// Config is an interface to manage particular helm release
type Config interface {
	Uniq() uniqname.UniqName
	HandleDependencies([]Config)
	Sync() (*release.Release, error)
	NotifySuccess()
	NotifyFailed()
	DryRun(bool)
	ChartDepsUpd() error
	In(a []Config) bool
	BuildValues(dir string, gomplate *template.GomplateConfig) error
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

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v2 to parse YAML into `Config` interface
func UnmarshalYAML(unmarshal func(interface{}) error) ([]Config, error) {
	r := make([]*config, 0)
	if err := unmarshal(&r); err != nil {
		return nil, err
	}

	res := make([]Config, len(r))
	for i := range r {
		res[i] = r[i]
	}

	return res, nil
}
