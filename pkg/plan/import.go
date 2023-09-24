package plan

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/version"
	log "github.com/sirupsen/logrus"
)

type PlanImportFS interface {
	fs.ReadDirFS
	fs.ReadFileFS
	fs.StatFS
	fs.SubFS
}

// Import parses directory with plan files and imports them into structure.
func (p *Plan) Import(ctx context.Context, baseFSUntyped fs.FS) error {
	baseFS, ok := baseFSUntyped.(PlanImportFS)
	if !ok {
		return fmt.Errorf("invalid plandir for import: %w", ErrInvalidPlandir)
	}

	yml, err := baseFS.Sub(File)
	if err != nil {
		return err
	}

	body, err := NewBody(ctx, yml, true)
	if err != nil {
		return err
	}

	err = p.importManifest(baseFS)

	switch {
	case errors.Is(err, ErrManifestDirEmpty), errors.Is(err, fs.ErrNotExist):
		log.WithError(err).Warn("error caught while importing manifests")
	case err != nil:
		return err
	}

	p.body = body

	// Validate all files exist.
	err = p.ValidateValuesImport(baseFS)
	if err != nil {
		return err
	}

	version.Validate(p.body.Version)

	return nil
}

func (p *Plan) importManifest(baseFS PlanImportFS) error {
	ls, err := baseFS.ReadDir(Manifest)
	if err != nil {
		return fmt.Errorf("failed to read manifest dir %s: %w", Manifest, err)
	}

	if len(ls) == 0 {
		return ErrManifestDirEmpty
	}

	for _, l := range ls {
		if l.IsDir() {
			continue
		}

		f := filepath.Join(Manifest, l.Name())
		c, err := baseFS.ReadFile(f)
		if err != nil {
			return fmt.Errorf("failed to read manifest %s: %w", f, err)
		}

		n := strings.TrimSuffix(l.Name(), filepath.Ext(l.Name())) // drop extension of file

		p.manifests[uniqname.UniqName(n)] = string(c)
	}

	return nil
}
