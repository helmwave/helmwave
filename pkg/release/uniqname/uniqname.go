package uniqname

import (
	"fmt"
	"regexp"
	"strings"
)

// Separator is a separator between release name and namespace.
const Separator = "@"

var validateRegexp = regexp.MustCompile("[a-z0-9](_[-a-z0-9]*[a-z0-9])?")

// UniqName is an alias for string.
type UniqName string

// Generate returns uniqname for provided release name and namespace.
func Generate(name, namespace, kubecontext string) (UniqName, error) {
	u := UniqName(fmt.Sprintf("%s%s%s%s", name, Separator, namespace, kubecontext))

	return u, u.Validate()
}

// Equal checks whether uniqnames are equal.
func (n UniqName) Equal(a UniqName) bool {
	return n == a
}

// Validate validates this object.
func (n UniqName) Validate() error {
	s := strings.Split(n.String(), Separator)
	if len(s) != 3 {
		return NewValidationError(n.String())
	}

	// I know, it should be just 3 items in slice
	for i := range s {
		if !validateRegexp.MatchString(s[i]) {
			return NewValidationError(n.String())
		}
	}

	return nil
}

func (n UniqName) String() string {
	return string(n)
}
