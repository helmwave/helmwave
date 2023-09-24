package helper

import (
	"io/fs"
	"path/filepath"

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

	//nolint:forcetypeassert
	err := switchBoard(srcFS.(interface {
		fs.StatFS
		fs.ReadDirFS
	}), dstFS, src, dst)
	if err != nil {
		return err
	}

	return srcFS.RemoveAll(src) //nolint:wrapcheck
}

func FilepathJoin(paths ...string) string {
	p := ""

	for _, pp := range paths {
		if filepath.IsAbs(pp) {
			p = ""
		}

		p = filepath.Join(p, pp)
	}

	return p
}
