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
	baseDir := filepath.Dir(p)
	if p == "" {
		baseDir = ".."
	}
	if err := f.MkdirAll(baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create base directories for %s: %w", p, err)
	}

	file, err := f.Create(p)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", p, err)
	}

	return file, nil
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

func CopyDir(srcFS interface {
	fs.StatFS
	fs.ReadDirFS
}, destFS fsimpl.WriteableFS, srcdir, destdir string) error {
	var contents []os.FileInfo
	entries, err := srcFS.ReadDir(srcdir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			return err
		}
		contents = append(contents, info)
	}

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, content := range contents {
		cs, cd := filepath.Join(srcdir, content.Name()), filepath.Join(destdir, content.Name())

		if err = switchBoard(srcFS, destFS, cs, cd); err != nil {
			// If any error, exit immediately
			return err
		}
	}

	return nil
}

func switchBoard(srcFS interface {
	fs.StatFS
	fs.ReadDirFS
}, destFS fsimpl.WriteableFS, srcpath, destpath string) error {
	stat, err := srcFS.Stat(srcpath)
	if err != nil {
		return err
	}

	switch {
	case stat.IsDir():
		return CopyDir(srcFS, destFS, srcpath, destpath)
	default:
		return CopyFile(srcFS, destFS, srcpath, destpath)
	}
}
