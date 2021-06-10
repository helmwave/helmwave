package plan

import "errors"

func (p *Plan) Rollback() error {
	return errors.New("sorry, rollback not ready yet")
}
