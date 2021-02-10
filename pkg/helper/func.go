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

func Save2File(p string, c string) error {
	f, err := CreateFile(p)
	if err != nil {
		return err
	}

	_, err = f.WriteString(c)
	if err != nil {
		return err
	}

	return f.Close()
}
