package helper

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func Contains(t string, a []string) bool {
	for _, v := range a {
		if v == t {
			return true
		}
	}

	return false
}

func CreateFile(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return nil, err
	}

	return os.Create(p)
}

// IsExists return true if file exists.
func IsExists(s string) bool {
	_, err := os.Stat(s)
	switch {
	case err == nil:
		return true
	case os.IsNotExist(err):
		return false
	default:
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		log.Fatal(err)

		return false
	}
}
