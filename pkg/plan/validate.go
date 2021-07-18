package plan

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
)

var ErrValidateFailed = errors.New("validate failed")

func (p *Plan) ValidateValues() error {
	f := false
	for _, rel := range p.body.Releases {
		for i := range rel.Values {
			_, err := os.Stat(rel.Values[i].Get())
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

	return ErrValidateFailed
}
