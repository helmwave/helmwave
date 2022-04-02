package plan

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/version"
	log "github.com/sirupsen/logrus"
)

// Import parses directory with plan files and imports them into structure.
func (p *Plan) Import() error {
	body, err := NewBody(p.fsys, p.File())
	if err != nil {
		return err
	}

	err = p.importManifest()
	if errors.Is(err, ErrManifestDirEmpty) {
		log.Warn(err)
	}

	if !errors.Is(err, ErrManifestDirEmpty) && err != nil {
		return err
	}

	p.body = body
	version.Check(p.body.Version, version.Version)

	return nil
}

func (p *Plan) importManifest() error {
	d := filepath.Join(p.URL.Path, Manifest)
	ls, err := fs.ReadDir(p.fsys, p.Dir())
	if err != nil {
		return fmt.Errorf("failed to read manifest dir %s: %w", d, err)
	}

	if len(ls) == 0 {
		return ErrManifestDirEmpty
	}

	for _, l := range ls {
		if l.IsDir() {
			continue
		}

		f := filepath.Join(p.Dir(), Manifest, l.Name())
		c, err := fs.ReadFile(p.fsys, f)
		if err != nil {
			return fmt.Errorf("failed to read manifest %s: %w", f, err)
		}

		n := strings.TrimSuffix(l.Name(), filepath.Ext(l.Name())) // drop extension of file

		p.manifests[uniqname.UniqName(n)] = string(c)
	}

	return nil
}
