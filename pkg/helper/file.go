package helper

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	cp "github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
)

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

func CreateWriteFile(path string, content string) (err error) {
	f, err := CreateFile(path)
	if err != nil {
		return err
	}

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return f.Close()
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
		// Schrödinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		log.Fatal(err)

		return false
	}
}

// CopyFile copy file to dest. Destination is either file or dir.
func CopyFile(src, dest string) error {
	destStat, err := os.Stat(dest)

	if err == nil {
		if destStat.Mode().IsDir() {
			dest = path.Join(dest, filepath.Base(src))
		} else {
			return fmt.Errorf("failed to copy file '%s': destination '%s' already exists", src, dest)
		}
	}

	err = cp.Copy(src, dest)
	if err != nil {
		return fmt.Errorf("failed to copy file '%s' to '%s': %w", src, dest, err)
	}

	return nil
}
