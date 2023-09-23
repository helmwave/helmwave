package plan

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
)

// Export allows save plan to file.
func (p *Plan) Export(ctx context.Context, baseFS fsimpl.WriteableFS, skipUnchanged bool) error {
	if err := baseFS.RemoveAll(p.dir); err != nil {
		return fmt.Errorf("failed to clean plan directory %s: %w", p.dir, err)
	}
	defer func(dir string) {
		err := baseFS.RemoveAll(dir)
		if err != nil {
			p.Logger().WithError(err).Error("failed to remove temporary directory")
		}
	}(p.tmpDir)

	if skipUnchanged {
		p.removeUnchanged()
		p.Logger().Info("removed unchanged releases from plan")
	}

	wg := parallel.NewWaitGroup()
	wg.Add(4)

	go func() {
		defer wg.Done()
		if err := p.exportCharts(baseFS); err != nil {
			wg.ErrChan() <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := p.exportManifest(baseFS); err != nil {
			wg.ErrChan() <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := p.exportValues(baseFS); err != nil {
			wg.ErrChan() <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := p.exportGraphMD(baseFS); err != nil {
			wg.ErrChan() <- err
		}
	}()

	err := wg.Wait()
	if err != nil {
		return err
	}

	// Save Planfile after everything is exported
	return helper.SaveInterface(ctx, baseFS, p.fullPath, p.body)
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

func (p *Plan) exportCharts(baseFS fsimpl.WriteableFS) error {
	for i, rel := range p.body.Releases {
		l := p.Logger().WithField("release", rel.Uniq())

		if !rel.Chart().IsRemote(baseFS) {
			l.Info("chart is local, skipping exporting it")

			continue
		}

		src := path.Join(p.tmpDir, "charts", rel.Uniq().String())
		dst := path.Join(p.dir, "charts", rel.Uniq().String())
		err := helper.MoveFile(
			baseFS,
			baseFS,
			src,
			dst,
		)
		if err != nil {
			return err
		}

		// Chart is places as an archive under this directory.
		// So we need to find it and use.
		entries, err := baseFS.(fs.ReadDirFS).ReadDir(dst)
		if err != nil {
			l.WithError(err).Warn("failed to read directory with downloaded chart, skipping")

			continue
		}

		if len(entries) != 1 {
			l.WithField("entries", entries).Warn("don't know which file is downloaded chart, skipping")

			continue
		}

		chart := entries[0]
		p.body.Releases[i].SetChartName(path.Join(dst, chart.Name()))
	}

	return nil
}

func (p *Plan) exportManifest(baseFS fsimpl.WriteableFS) error {
	for k, v := range p.manifests {
		m := filepath.Join(p.dir, Manifest, k.String()+".yml")

		f, err := helper.CreateFile(baseFS, m)
		if err != nil {
			return err
		}

		_, err = f.Write([]byte(v))
		if err != nil {
			return fmt.Errorf("failed to write manifest %s: %w", m, err)
		}

		err = f.Close()
		if err != nil {
			return fmt.Errorf("failed to close manifest %s: %w", m, err)
		}
	}

	return nil
}

func (p *Plan) exportGraphMD(baseFS fsimpl.WriteableFS) error {
	found := false
	for _, rel := range p.body.Releases {
		if len(rel.DependsOn()) > 0 {
			found = true

			break
		}
	}

	if !found {
		return nil
	}

	const filename = "graph.md"
	filePath := filepath.Join(p.dir, filename)
	f, err := helper.CreateFile(baseFS, filePath)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(p.graphMD))
	if err != nil {
		return fmt.Errorf("failed to write graph file %s: %w", filePath, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close graph file %s: %w", filePath, err)
	}

	return nil
}

func (p *Plan) exportValues(baseFS fsimpl.WriteableFS) error {
	found := false

	for i, rel := range p.body.Releases {
		for j := range p.body.Releases[i].Values() {
			found = true
			p.body.Releases[i].Values()[j].SetUniq(p.dir, rel.Uniq())
		}
	}

	if !found {
		return nil
	}

	// It doesnt work if workdir has been mounted.
	err := helper.MoveFile(
		baseFS,
		baseFS,
		filepath.Join(p.tmpDir, Values),
		filepath.Join(p.dir, Values),
	)
	if err != nil {
		return fmt.Errorf("failed to copy values from %s to %s: %w", p.tmpDir, p.dir, err)
	}

	return nil
}

// IsExist returns true if planfile exists.
func (p *Plan) IsExist(baseFS fs.StatFS) bool {
	return helper.IsExists(baseFS, p.fullPath)
}

// IsManifestExist returns true if planfile exists.
func (p *Plan) IsManifestExist(baseFS fs.StatFS) bool {
	return helper.IsExists(baseFS, filepath.Join(p.dir, Manifest))
}
