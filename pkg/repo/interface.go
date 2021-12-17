package repo

import (
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"
)

type Config interface {
	In([]Config) bool
	Install(*helm.EnvSettings, *repo.File) error
	Name() string
	URL() string
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
