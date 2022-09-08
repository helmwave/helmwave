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
	r := new(jsonschema.Reflector)
	var l []*Registry
	return r.Reflect(&l)
}

// Registry is main registry.
type Registry struct {
	log      *log.Entry
	HostF    string `json:"host" jsonschema:"required"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Insecure bool   `json:"insecure,omitempty"`
}

// Host return Host value.
func (c *Registry) Host() string {
	return c.HostF
}

// func (c *Registry) Username() string {
//	return c.UsernameF
// }
//
// func (c *Registry) Password() string {
//	return c.PasswordF
// }

func (c *Registry) Logger() *log.Entry {
	if c.log == nil {
		c.log = log.WithField("registry", c.Host())
	}

	return c.log
}
