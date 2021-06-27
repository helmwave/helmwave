package helper

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sort"
)

func Contains(t string, a []string) bool {
	i := sort.SearchStrings(a, t)
	return i < len(a) && a[i] == t
}

func CreateFile(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return nil, err
	}
	return os.Create(p)
}

// Inclusion checks where any of release tags are included in target tags.
func Inclusion(where, that []string, matchAll bool) bool {
	if len(where) == 0 {
		return true
	}

	for _, t := range where {
		contains := Contains(t, that)
		if matchAll && !contains {
			return false
		}
		if !matchAll && contains {
			return true
		}
	}

	return matchAll
}

func IsExists(s string) bool {
	if _, err := os.Stat(s); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		log.Fatal(err)
		return false
	}

}
