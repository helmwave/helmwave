package plan

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
)

// Export allows save plan to file.
func (p *Plan) Export(ctx context.Context, srcFS fsimpl.CurrentPathFS, plandirFSUntyped fs.FS, skipUnchanged bool) error {
	plandirFS, ok := plandirFSUntyped.(ExportFS)
	if !ok {
		return fmt.Errorf("invalid plandir for export: %w", ErrInvalidPlandir)
	}

	if err := plandirFS.RemoveAll(""); err != nil {
		return fmt.Errorf("failed to clean plan directory: %w", err)
	}

	if skipUnchanged {
		p.removeUnchanged()
		p.Logger().Info("removed unchanged releases from plan")
	}

	wg := parallel.NewWaitGroup()
	wg.Add(3)

	go func() {
		defer wg.Done()
		if err := p.exportCharts(srcFS, plandirFS); err != nil {
			wg.ErrChan() <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := p.exportValues(srcFS, plandirFS); err != nil {
			wg.ErrChan() <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := p.exportGraphMD(plandirFS); err != nil {
			wg.ErrChan() <- err
		}
	}()

	err := wg.Wait()
	if err != nil {
		return err
	}

	err = p.exportManifests(ctx, plandirFS)
	if err != nil {
		return err
	}

	// Save Planfile after everything is exported
	return helper.SaveInterface(ctx, plandirFS, p.fullPath, p.body)
}

func (p *Plan) removeUnchanged() {
	filtered := p.body.Releases[:0]

	for _, rel := range p.body.Releases {
		if !helper.In[release.Config](rel, p.unchanged) {
			filtered = append(filtered, rel)
		}
	}

	p.body.Releases = filtered
}

// IsExist returns true if planfile exists.
func (p *Plan) IsExist(baseFS fs.FS) bool {
	return helper.IsExists(baseFS, p.fullPath)
}

// IsManifestExist returns true if planfile exists.
func (p *Plan) IsManifestExist(baseFS fs.FS) bool {
	return helper.IsExists(baseFS, Manifest)
}
