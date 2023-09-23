package helper

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/helmwave/go-fsimpl"
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
func CreateFile(f fsimpl.WriteableFS, p string) (fsimpl.WriteableFile, error) {
	if err := f.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return nil, fmt.Errorf("failed to create base directories for %s: %w", p, err)
	}

	file, err := f.Create(p)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", p, err)
	}

	return file, nil
}

func Path(f fs.FS) (string, error) {
	file, err := f.Open(".")
	if err != nil {
		return "", err //nolint:wrapcheck
	}
	defer func() {
		_ = file.Close()
	}()

	stat, err := file.Stat()
	if err != nil {
		return "", err //nolint:wrapcheck
	}

	return stat.Name(), nil
}

// IsExists return true if file exists.
func IsExists(f fs.StatFS, s string) bool {
	_, err := f.Stat(s)
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
func CopyFile(srcFS fs.FS, destFS fsimpl.WriteableFS, src, dest string) error {
	destStat, err := destFS.Stat(dest)

	if err == nil {
		if destStat.Mode().IsDir() {
			dest = path.Join(dest, filepath.Base(src))
		} else {
			return fmt.Errorf("failed to copy file '%s': destination '%s' already exists", src, dest)
		}
	}

	readcloser, err := srcFS.Open(src)
	if err != nil {
		return fmt.Errorf("failed to copy file '%s': %w", src, err)
	}
	defer func() {
		err := readcloser.Close()
		if err != nil {
			log.WithError(err).Error("failed to close file")
		}
	}()

	err = destFS.MkdirAll(filepath.Dir(dest), fs.ModePerm)
	if err != nil {
		return err //nolint:wrapcheck
	}

	f, err := destFS.Create(dest)
	if err != nil {
		return err //nolint:wrapcheck
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.WithError(err).Error("failed to close file")
		}
	}()

	var buf []byte = nil
	var w io.Writer = f
	var r io.Reader = readcloser

	_, err = io.CopyBuffer(w, r, buf)
	if err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}
