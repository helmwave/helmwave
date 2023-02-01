package helper

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Contains checks whether string exists in string slice.
func Contains(t string, a []string) bool {
	for _, v := range a {
		if v == t {
			return true
		}
	}

	return false
}

// CreateFile creates recursively basedir of file and returns created file object.
func CreateFile(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return nil, fmt.Errorf("failed to create base directories for %s: %w", p, err)
	}

	f, err := os.Create(p)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", p, err)
	}

	return f, nil
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

// CopyFile copy file to dest. Destination is either file or dir.
func CopyFile(src, dest string) error {
	fin, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", src, err)
	}
	defer func() {
		e := fin.Close()
		// Report the error from Close, if any.
		// But do so only if there isn't already
		// an outgoing error.
		if e != nil && err == nil {
			err = e
		}
	}()

	destStat, err := os.Stat(dest)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", dest, err)
	}

	// if destination is directory, add file name
	if destStat.Mode().IsDir() {
		dest = fmt.Sprintf("%s/%s", dest, filepath.Base(src))
	}

	fout, err := CreateFile(dest)
	if err != nil {
		return err
	}
	defer func() {
		e := fout.Close()
		if e != nil && err == nil {
			err = e
		}
	}()

	_, copyErr := io.Copy(fout, fin)
	if copyErr != nil {
		return fmt.Errorf("failed to copy file %s: %w", src, copyErr)
	}

	return err
}
