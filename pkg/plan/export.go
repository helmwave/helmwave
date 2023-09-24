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
	log "github.com/sirupsen/logrus"
)

// Export allows save plan to file.
func (p *Plan) Export(ctx context.Context, plandirFSUntyped fs.FS, skipUnchanged bool) error {
	plandirFS, ok := plandirFSUntyped.(fsimpl.WriteableFS)
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
	wg.Add(4)

	go func() {
		defer wg.Done()
		if err := p.exportCharts(plandirFS); err != nil {
			wg.ErrChan() <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := p.exportManifest(plandirFS); err != nil {
			wg.ErrChan() <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := p.exportValues(plandirFS); err != nil {
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

func (p *Plan) exportCharts(baseFS fsimpl.WriteableFS) error {
	for i, rel := range p.body.Releases {
		l := p.Logger().WithField("release", rel.Uniq())

		if !rel.Chart().IsRemote(baseFS) {
			l.Info("chart is local, skipping exporting it")

			continue
		}

		dst := path.Join(Charts, rel.Uniq().String())
		err := rel.DownloadChart(baseFS, baseFS, dst)
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
		m := filepath.Join(Manifest, k.String()+".yml")

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
	f, err := helper.CreateFile(baseFS, filename)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(p.graphMD))
	if err != nil {
		return fmt.Errorf("failed to write graph file %s: %w", filename, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close graph file %s: %w", filename, err)
	}

	return nil
}

func (p *Plan) exportValues(plandirFS fsimpl.WriteableFS) error {
	//found := false

	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for _, rel := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel release.Config) {
			defer wg.Done()
			err := rel.ExportValues(plandirFS, plandirFS, p.templater)
			if err != nil {
				log.Errorf("❌ %s values: %v", rel.Uniq(), err)
				wg.ErrChan() <- err
			} else {
				var vals []string
				for i := range rel.Values() {
					vals = append(vals, rel.Values()[i].Src)
				}

				if len(vals) == 0 {
					rel.Logger().Info("🔨 no values provided")
				} else {
					rel.Logger().WithField("values", vals).Infof("✅ found %d values count", len(vals))
				}
			}
		}(wg, rel)
		//for j := range p.body.Releases[i].Values() {
		//	found = true
		//	p.body.Releases[i].Values()[j].getUniqPath(rel.Uniq())
		//}
	}

	return wg.Wait()

	//if !found {
	//	return nil
	//}
	//
	//// It doesnt work if workdir has been mounted.
	//err := helper.MoveFile(
	//	plandirFS,
	//	plandirFS,
	//	filepath.Join(p.tmpDir, Values),
	//	filepath.Join(Values),
	//)
	//if err != nil {
	//	return fmt.Errorf("failed to copy values from %s: %w", p.tmpDir, err)
	//}

	//return nil
}

// IsExist returns true if planfile exists.
func (p *Plan) IsExist(baseFS fs.StatFS) bool {
	return helper.IsExists(baseFS, p.fullPath)
}

// IsManifestExist returns true if planfile exists.
func (p *Plan) IsManifestExist(baseFS fs.StatFS) bool {
	return helper.IsExists(baseFS, Manifest)
}
