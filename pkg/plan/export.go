package plan

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

// Export allows save plan to file.
func (p *Plan) Export(ctx context.Context, skipUnchanged bool) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	log.Tracef("I am exporting plan to %s", filepath.Join(wd, p.dir))

	if err := os.RemoveAll(p.dir); err != nil {
		return fmt.Errorf("failed to clean plan directory %s: %w", p.dir, err)
	}
	defer func(dir string) {
		err := os.RemoveAll(dir)
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
	}()
	go func() {
		defer wg.Done()
		if err := p.exportGraphMD(); err != nil {
			wg.ErrChan() <- err
		}
	}()

	err = wg.Wait()
	if err != nil {
		return err
	}

	// Save Planfile after everything is exported
	return helper.SaveInterface(ctx, p.fullPath, p.body)
}

func (p *Plan) removeUnchanged() {
	p.body.Releases = slices.DeleteFunc(p.body.Releases, func(rel release.Config) bool {
		return slices.ContainsFunc(p.unchanged, func(r release.Config) bool {
			return r.Uniq().Equal(rel.Uniq())
		})
	})
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
		p.body.Releases[i].SetChartName(path.Join(dst, chart.Name()))
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
	found := slices.ContainsFunc(p.body.Releases, func(rel release.Config) bool {
		return len(rel.DependsOn()) > 0
	})
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

	for i := range p.body.Releases {
		for j := range p.body.Releases[i].Values() {
			found = true
			v := p.body.Releases[i].Values()[j]
			//nolint: govet
			v.Dst = filepath.Join(p.dir, "values", p.body.Releases[i].Uniq().String(), strconv.Itoa(i)+".yml")
		}
	}

	if !found {
		return nil
	}

	// It doesn't work if workdir has been mounted.
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
