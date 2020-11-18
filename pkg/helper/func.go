package helper

import (
	"sort"
	"unicode/utf8"
)

func Contains(t string, a []string) bool {
	i := sort.SearchStrings(a, t)
	return i < len(a) && a[i] == t
}

func TrimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
