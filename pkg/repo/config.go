package repo

import (
	"errors"

	"helm.sh/helm/v3/pkg/repo"

	log "github.com/sirupsen/logrus"
)

// ErrNotFound is an error for not declared repository name.
var ErrNotFound = errors.New("repository not found")

type config struct {
	log        *log.Entry `yaml:"-" json:"-"`
	repo.Entry `yaml:",inline" json:",inline"`
	Force      bool `yaml:"force" json:"force" jsonschema:"title=force flag,description=force update helm repo list and download dependencies,default=false"`
}

func (c *config) Name() string {
	return c.Entry.Name
}

func (c *config) URL() string {
	return c.Entry.URL
}

func (c *config) Logger() *log.Entry {
	if c.log == nil {
		c.log = log.WithField("repository", c.Name())
	}

	return c.log
}
