package uniqname

import (
	"fmt"
	"regexp"
	"strings"
)

// Separator is a separator between release name and namespace.
const Separator = "@"

var validateRegexp = regexp.MustCompile("[a-z0-9]([-a-z0-9]*[a-z0-9])?")

// UniqName is an alias for string.
type UniqName string

// Generate returns uniqname for provided release name and namespace.
func Generate(name, namespace string) (UniqName, error) {
	u := UniqName(fmt.Sprintf("%s%s%s", name, Separator, namespace))

	return u, u.Validate()
}

// GenerateWithDefaultNamespace parses uniqname out of provided line.
// If there is no namespace in line, default namespace will be used.
func GenerateWithDefaultNamespace(line, namespace string) (UniqName, error) {
	s := strings.Split(line, Separator)

	name := s[0]

	if len(s) > 1 && s[1] != "" {
		namespace = s[1]
	}

	return Generate(name, namespace)
}

// Equal checks whether uniqnames are equal.
func (n UniqName) Equal(a UniqName) bool {
	return n == a
}

// Validate validates this object.
func (n UniqName) Validate() error {
	s := strings.Split(n.String(), Separator)
	if len(s) != 2 {
		return NewValidationError(n.String())
	}

	if !validateRegexp.MatchString(s[0]) {
		return NewValidationError(n.String())
	}

	if !validateRegexp.MatchString(s[1]) {
		return NewValidationError(n.String())
	}

	return nil
}

func (n UniqName) String() string {
	return string(n)
}
