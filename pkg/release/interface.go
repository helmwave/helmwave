package release

import (
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/template"
	"helm.sh/helm/v3/pkg/release"
)

type Config interface {
	Uniq() uniqname.UniqName
	HandleDependencies([]Config)
	Sync() (*release.Release, error)
	NotifySuccess()
	NotifyFailed()
	DryRun(bool) Config
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
