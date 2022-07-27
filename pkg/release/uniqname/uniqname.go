package uniqname

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Separator is a separator between release name and namespace.
const Separator = "@"

// ErrValidate is an error for failed uniqname validation.
var ErrValidate = errors.New("failed to validate uniqname")

// UniqName is an alias for string.
type UniqName string

// Generate returns uniqname for provided release name and namespace.
func Generate(name, namespace string) (UniqName, error) {
	u := UniqName(fmt.Sprintf("%s%s%s", name, Separator, namespace))

	return u, u.Validate()
}

// Contains searches for uniqname in slice of uniqnames.
func Contains(t UniqName, a []UniqName) bool {
	for _, v := range a {
		if v == t {
			return true
		}
	}

	return false
}

// Equal checks whether uniqnames are equal.
func (n UniqName) Equal(a UniqName) bool {
	return n == a
}

// Validate validates this object.
func (n UniqName) Validate() error {
	s := string(n)
	if len(strings.Split(s, Separator)) != 2 {
		return ErrValidate
	}

	r := regexp.MustCompile("[a-z0-9]([-a-z0-9]*[a-z0-9])?" + Separator + "[a-z0-9]([-a-z0-9]*[a-z0-9])?")

	if !r.MatchString(s) {
		return ErrValidate
	}

	return nil
}
