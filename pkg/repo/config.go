package repo

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/repo"
)

// Configs type of array Config.
type Configs []Config

// UnmarshalYAML parse Config.
func (r *Configs) UnmarshalYAML(node *yaml.Node) error {
	if r == nil {
		r = new(Configs)
	}
	var err error

	*r, err = UnmarshalYAML(node)

	return err
}

type config struct {
	log        *log.Entry       `yaml:"-"`
	repo.Entry `yaml:",inline"` //nolint:nolintlint
	Force      bool             `yaml:"force"`
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

// ErrNotFound is an error for not declared repository name.
var ErrNotFound = errors.New("repository not found")
