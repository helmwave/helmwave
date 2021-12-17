package repo

import (
	"errors"

	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"
)

type Config interface {
	In([]Config) bool
	Install(*helm.EnvSettings, *repo.File) error
	Name() string
	URL() string
}

type config struct {
	repo.Entry `yaml:",inline"`
	Force      bool
}

func (c *config) Name() string {
	return c.Entry.Name
}

func (c *config) URL() string {
	return c.Entry.URL
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

var ErrNotFound = errors.New("repository not found")
