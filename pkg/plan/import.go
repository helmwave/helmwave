package plan

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/version"
	log "github.com/sirupsen/logrus"
)

// Import parses directory with plan files and imports them into structure.
func (p *Plan) Import(ctx context.Context) error {
	body, err := NewBody(ctx, p.fullPath, true)
	if err != nil {
		return err
	}

	err = p.importManifest()

	switch {
	case errors.Is(err, ErrManifestDirEmpty), errors.Is(err, fs.ErrNotExist):
		log.WithError(err).Warn("error caught while importing manifests")
	case err != nil:
		return err
	}

	p.body = body

	// Validate all files exist.
	err = p.ValidateValuesImport()
	if err != nil {
		return err
	}

	version.Validate(p.body.Version)

	return nil
}

func (p *Plan) importManifest() error {
	d := filepath.Join(p.dir, Manifest)
	ls, err := os.ReadDir(d)
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

		f := filepath.Join(p.dir, Manifest, l.Name())
		c, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("failed to read manifest %s: %w", f, err)
		}

		n := strings.TrimSuffix(l.Name(), filepath.Ext(l.Name())) // drop extension of file

		u, _ := uniqname.NewFromString(n) // we don't expect errors here
		p.manifests[u] = string(c)
	}

	return nil
}
