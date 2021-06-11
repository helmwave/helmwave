package plan

import (
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	log "github.com/sirupsen/logrus"
	"os"
)

const planfile = "planfile"

type Plan struct {
	body     *planBody
	dir      string
	fullPath string
}

type planBody struct {
	Project      string
	Version      string
	Repositories []*repo.Config
	Releases     []*release.Config
}

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

func New(dir string) *Plan {
	if dir[len(dir)-1:] != "/" {
		dir += "/"
	}

	plan := &Plan{
		dir:      dir,
		fullPath: dir + planfile,
	}

	return plan
}
