package registry

import (
	"github.com/invopop/jsonschema"
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
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
	}
	var l []*config

	return r.Reflect(&l)
}

// config is main registry config.
//
//nolint:lll
type config struct {
	log      *log.Entry `json:"-"`
	HostF    string     `json:"host" jsonschema:"required,description=OCI registry host optionally with port,pattern=^.*(:[0-9]+)?$"`
	Username string     `json:"username"`
	Password string     `json:"password"`
	Insecure bool       `json:"insecure" jsonschema:"default=false"`
}

// Host return Host value.
func (c *config) Host() string {
	return c.HostF
}

func (c *config) Logger() *log.Entry {
	if c.log == nil {
		c.log = log.WithField("registry", c.Host())
	}

	return c.log
}
