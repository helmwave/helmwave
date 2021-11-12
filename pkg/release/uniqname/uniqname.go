package uniqname

import (
	"regexp"
	"strings"
)

const Separator = "@"

type UniqName string

func Contains(t UniqName, a []UniqName) bool {
	for _, v := range a {
		if v == t {
			return true
		}
	}

	return false
}

func (n UniqName) In(a []UniqName) bool {
	for _, v := range a {
		if v == n {
			return true
		}
	}

	return false
}

func (n UniqName) Validate() bool {
	s := string(n)
	if len(strings.Split(s, Separator)) != 2 {
		return false
	}

	r := regexp.MustCompile("[a-z0-9]([-a-z0-9]*[a-z0-9])?" + Separator + "[a-z0-9]([-a-z0-9]*[a-z0-9])?")

	return r.MatchString(s)
}
