package action

import (
	"errors"
	"flag"
	"io/fs"
	"net/url"
	"path/filepath"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/blobfs"
	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/helmwave/go-fsimpl/gitfs"
	"github.com/urfave/cli/v2"
)

var (
	_ cli.Generic = (*genericFS[fs.FS])(nil)
	_ flag.Getter = (*genericFS[fs.FS])(nil)

	ErrInvalidFilesystem = errors.New("this filesystem is not supported for this flag")

	mux = fsimpl.NewMux()
)

func init() {
	mux.Add(filefs.FS)
	mux.Add(blobfs.FS)
	mux.Add(gitfs.FS)
}

type genericFS[T fs.FS] struct {
	dest *T
	val  string
}

func (f *genericFS[T]) String() string {
	return f.val
}

func (f *genericFS[T]) Set(val string) error {
	f.val = val

	resFS, err := mux.Lookup(stringFSToURLFS(val).String())
	if err != nil {
		return err
	}

	castedFS, ok := resFS.(T)
	if !ok {
		return ErrInvalidFilesystem
	}

	*f.dest = castedFS
	return nil
}

func (f *genericFS[T]) Get() any {
	return *f.dest
}

func stringFSToURLFS(val string) *url.URL {
	u, err := url.Parse(val)
	if err == nil && u.Scheme != "" {
		return u
	}

	absPath, _ := filepath.Abs(val)
	return &url.URL{Scheme: "file", Path: absPath}
}

func createGenericFS[T fs.FS](dest *T, subPaths ...string) *genericFS[T] {
	f := &genericFS[T]{
		dest: dest,
	}
	_ = f.Set(getDefaultFSValue(subPaths...))

	return f
}

func getDefaultFSValue(subPaths ...string) string {
	if len(subPaths) == 0 {
		return "."
	}

	return filepath.Join(subPaths...)
}
