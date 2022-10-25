package repo

import (
	"errors"

	"helm.sh/helm/v3/pkg/repo"

	"github.com/invopop/jsonschema"
	log "github.com/sirupsen/logrus"
)

func (Configs) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
	}
	var l []*config

	return r.Reflect(&l)
}

//nolint:lll
type config struct {
	log        *log.Entry `json:"-"`
	repo.Entry `json:",inline"`
	Force      bool `json:"force" jsonschema:"title=force flag,description=force update helm repo list and download dependencies,default=false"`
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
