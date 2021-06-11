package helper

import (
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

// CheckTagInclusion checks where any of release tags are included in target tags.
func CheckTagInclusion(targetTags, releaseTags []string, matchAll bool) bool {
	if len(targetTags) == 0 {
		return true
	}

	for _, t := range targetTags {
		contains := Contains(t, releaseTags)
		if matchAll && !contains {
			return false
		}
		if !matchAll && contains {
			return true
		}
	}

	return matchAll
}
