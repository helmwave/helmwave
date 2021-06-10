package helmwave

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func Command404(c *cli.Context, s string) {
	log.Errorf("👻 Command %q not found \n", s)
	os.Exit(127)
}
