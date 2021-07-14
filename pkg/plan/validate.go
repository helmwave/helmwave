package plan

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
)

var ErrValidateFailed = errors.New("validate failed")

func (p *Plan) ValidateValues() error {
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

	return ErrValidateFailed
}
