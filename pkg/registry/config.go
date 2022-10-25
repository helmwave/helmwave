package registry

import (
	log "github.com/sirupsen/logrus"
)

// config is main registry config.
type config struct {
	log      *log.Entry `json:"-"`
	HostF    string     `json:"host" jsonschema:"required"`
	Username string     `json:"username"`
	Password string     `json:"password"`
	Insecure bool       `json:"insecure"`
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
