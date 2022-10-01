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
	r := &jsonschema.Reflector{DoNotReference: true}
	var l []*config

	return r.Reflect(&l)
}

// config is main registry config.
type config struct {
	log      *log.Entry `yaml:"-"`
	HostF    string     `yaml:"host" json:"host" jsonschema:"required"`
	Username string     `yaml:"username" json:"username"`
	Password string     `yaml:"password" json:"password"`
	Insecure bool       `yaml:"insecure" json:"insecure"`
}

// Host return Host value.
func (c *config) Host() string {
	return c.HostF
}

// func (c *config) Username() string {
//	return c.UsernameF
// }
//
// func (c *config) Password() string {
//	return c.PasswordF
// }

func (c *config) Logger() *log.Entry {
	if c.log == nil {
		c.log = log.WithField("registry", c.Host())
	}

	return c.log
}
