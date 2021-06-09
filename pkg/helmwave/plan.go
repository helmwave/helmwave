package helmwave

import (
	log "github.com/sirupsen/logrus"
)

const PLANFILE = "planfile"
const MANIFEST = ".manifest/"

func (c *Config) InitPlan() {
	if c.Plandir[len(c.Plandir)-1:] != "/" {
		c.Plandir += "/"
	}
	log.Info("🛠 Your planfile is ", c.Plandir+PLANFILE)
}
