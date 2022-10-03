package repo

import (
	"errors"

	"github.com/invopop/jsonschema"
	"helm.sh/helm/v3/pkg/repo"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Configs type of array Config.
type Configs []Config

// UnmarshalYAML parse Config.
func (r *Configs) UnmarshalYAML(node *yaml.Node) error {
	var err error
	*r, err = UnmarshalYAML(node)

	return err
}

func (Configs) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{DoNotReference: true}
	var l []*config

	return r.Reflect(&l)
}

//nolint:lll
type config struct {
	log        *log.Entry `yaml:"-"`
	repo.Entry `yaml:",inline"`
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

// ErrNotFound is an error for not declared repository name.
var ErrNotFound = errors.New("repository not found")
