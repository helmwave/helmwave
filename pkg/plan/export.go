package plan

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
)

// Export allows save plan to file.
func (p *Plan) Export() error {
	if err := os.RemoveAll(p.dir); err != nil {
		return fmt.Errorf("failed to clean plan directory %s: %w", p.dir, err)
	}
	defer func(dir string) {
		err := os.RemoveAll(dir)
		if err != nil {
			p.Logger().WithError(err).Error("failed to remove temporary directory")
		}
	}(p.tmpDir)

	wg := parallel.NewWaitGroup()
	wg.Add(4)

	go func() {
		defer wg.Done()
		if err := p.exportCharts(); err != nil {
			wg.ErrChan() <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := p.exportManifest(); err != nil {
			wg.ErrChan() <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := p.exportValues(); err != nil {
			wg.ErrChan() <- err
		}

		// Save Planfile after values
		if err := helper.SaveInterface(p.fullPath, p.body); err != nil {
			wg.ErrChan() <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := p.exportGraphMD(); err != nil {
			wg.ErrChan() <- err
		}
	}()

	return wg.Wait()
}

func (p *Plan) exportCharts() error {
	for i, rel := range p.body.Releases {
		l := p.Logger().WithField("release", rel.Uniq())

		if !rel.Chart().IsRemote() {
			l.Info("chart is local, skipping exporting it")

			continue
		}

		src := path.Join(p.tmpDir, "charts", rel.Uniq().String())
		dst := path.Join(p.dir, "charts", rel.Uniq().String())
		err := helper.MoveFile(
			src,
			dst,
		)
		if err != nil {
			return err
		}

		// Chart is places as an archive under this directory.
		// So we need to find it and use.
		entries, err := os.ReadDir(dst)
		if err != nil {
			l.WithError(err).Warn("failed to read directory with downloaded chart, skipping")

			continue
		}

		if len(entries) != 1 {
			l.WithField("entries", entries).Warn("don't know which file is downloaded chart, skipping")

			continue
		}

		chart := entries[0]
		p.body.Releases[i].SetChart(path.Join(dst, chart.Name()))
	}

	return nil
}

func (p *Plan) exportManifest() error {
	for k, v := range p.manifests {
		m := filepath.Join(p.dir, Manifest, k.String()+".yml")

		f, err := helper.CreateFile(m)
		if err != nil {
			return err
		}

		_, err = f.WriteString(v)
		if err != nil {
			return fmt.Errorf("failed to write manifest %s: %w", f.Name(), err)
		}

		err = f.Close()
		if err != nil {
			return fmt.Errorf("failed to close manifest %s: %w", f.Name(), err)
		}
	}

	return nil
}

func (p *Plan) exportGraphMD() error {
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
	f, err := helper.CreateFile(filepath.Join(p.dir, filename))
	if err != nil {
		return err
	}

	_, err = f.WriteString(p.graphMD)
	if err != nil {
		return fmt.Errorf("failed to write graph file %s: %w", f.Name(), err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close graph file %s: %w", f.Name(), err)
	}

	return nil
}

func (p *Plan) exportValues() error {
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
		filepath.Join(p.tmpDir, Values),
		filepath.Join(p.dir, Values),
	)
	if err != nil {
		return fmt.Errorf("failed to copy values from %s to %s: %w", p.tmpDir, p.dir, err)
	}

	return nil
}

// IsExist returns true if planfile exists.
func (p *Plan) IsExist() bool {
	return helper.IsExists(p.fullPath)
}

// IsManifestExist returns true if planfile exists.
func (p *Plan) IsManifestExist() bool {
	return helper.IsExists(filepath.Join(p.dir, Manifest))
}
