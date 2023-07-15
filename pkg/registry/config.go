package registry

import (
	log "github.com/sirupsen/logrus"
)

// config is main registry config.
//
//nolint:lll
type config struct {
	log      *log.Entry `yaml:"-" json:"-"`
	HostF    string     `yaml:"host" json:"host" jsonschema:"required,description=OCI registry host optionally with port,pattern=^.*(:[0-9]+)?$"`
	Username string     `yaml:"username" json:"username"`
	Password string     `yaml:"password" json:"password"`
	Insecure bool       `yaml:"insecure" json:"insecure" jsonschema:"default=false"`
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
