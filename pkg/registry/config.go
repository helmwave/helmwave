package registry

import (
	log "github.com/sirupsen/logrus"
)

type config struct {
	log      *log.Entry `yaml:"-"`
	HostF    string     `yaml:"host"`
	Username string     `yaml:"username"`
	Password string     `yaml:"password"`
	Insecure bool       `yaml:"insecure"`
}

func (c *config) Host() string {
	return c.HostF
}

func (c *config) Logger() *log.Entry {
	if c.log == nil {
		c.log = log.WithField("registry", c.Host)
	}

	return c.log
}
