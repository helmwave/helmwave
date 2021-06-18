package plan

import (
	"errors"
	"github.com/helmwave/helmwave/pkg/version"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func (p *Plan) Import() error {
	src, err := ioutil.ReadFile(p.fullPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(src, p.body)
	if err != nil {
		return err
	}

	version.Check(p.body.Version, version.Version)

	return nil
}

// Validate that files is existing
func (p *Plan) Validate() error {
	f := false
	for _, rel := range p.body.Releases {
		for _, val := range rel.Values {
			_, err := os.Stat(val)
			if os.IsNotExist(err) {
				log.Error(err)
				f = true
			} else {
				// FatalError
				return err
			}
		}
	}
	if !f {
		return nil
	}

	return errors.New("cannot validate")
}
