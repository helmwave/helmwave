package uniqname

import (
	"github.com/pkg/errors"
)

func (n UniqName) Error(part string) error {
	return errors.Errorf("failed to validate uniqname: %s, problem with: %s", n.String(), part)
}
