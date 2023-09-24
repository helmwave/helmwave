package helper

import (
	"io/fs"

	"github.com/helmwave/go-fsimpl"
)

// MoveFile moves files or directories. It also handles move between different mounts (copy + rm).
func MoveFile(srcFS, dstFS fsimpl.WriteableFS, src, dst string) error {
	if srcFS == dstFS {
		err := dstFS.Rename(src, dst)
		if err == nil {
			return nil
		}
	}

	err := switchBoard(srcFS.(interface {
		fs.StatFS
		fs.ReadDirFS
	}), dstFS, src, dst)
	if err != nil {
		return err
	}

	return srcFS.RemoveAll(src)
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
