package helper

import (
	"os"
	"path/filepath"
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

func CreateFile(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
