package uniqname

import (
	"regexp"
	"strings"
)

// Separator is a separator between release name and namespace.
const Separator = "@"

// UniqName is an alias for string.
type UniqName string

// Contains searches for uniqname in slice of uniqnames.
func Contains(t UniqName, a []UniqName) bool {
	for _, v := range a {
		if v == t {
			return true
		}
	}

	return false
}

// In searches for uniqname in slice of uniqnames.
func (n UniqName) In(a []UniqName) bool {
	for _, v := range a {
		if v == n {
			return true
		}
	}

	return false
}

// Validate validates this object.
func (n UniqName) Validate() bool {
	s := string(n)
	if len(strings.Split(s, Separator)) != 2 {
		return false
	}

	r := regexp.MustCompile("[a-z0-9]([-a-z0-9]*[a-z0-9])?" + Separator + "[a-z0-9]([-a-z0-9]*[a-z0-9])?")

	return r.MatchString(s)
}
