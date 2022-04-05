package registry

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Configs []Config

func (r *Configs) UnmarshalYAML(node *yaml.Node) error {
	if r == nil {
		r = new(Configs)
	}
	var err error

	*r, err = UnmarshalYAML(node)

	return err
}

// config is main registry config
type config struct {
	log      *log.Entry `yaml:"-"`
	HostF    string     `yaml:"host"`
	Username string     `yaml:"username"`
	Password string     `yaml:"password"`
	Insecure bool       `yaml:"insecure"`
}

// Host return Host value
func (c *config) Host() string {
	return c.HostF
}

//func (c *config) Username() string {
//	return c.UsernameF
//}
//
//func (c *config) Password() string {
//	return c.PasswordF
//}

func (c *config) Logger() *log.Entry {
	if c.log == nil {
		c.log = log.WithField("registry", c.Host())
	}

	return c.log
}
