package helper

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

func Download(file, url string) error {
	f, err := CreateFile(file)
	if err != nil {
		return err
	}

	r, err := http.Get(url)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", r.Status)
	}

	_, err = io.Copy(f, r.Body)
	if err != nil {
		return err
	}

	return nil
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
