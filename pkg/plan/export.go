package plan

import (
	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"os"
)

// Export allows save plan to file
func (p *Plan) Export() error {
	if err := os.RemoveAll(p.dir); err != nil {
		return err
	}
	return helper.Save(p.fullPath, p.body)
}

// IsExist returns true if planfile exists
func (p *Plan) IsExist() bool {
	if _, err := os.Stat(p.fullPath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		log.Fatal(err)
		return false
	}
}
