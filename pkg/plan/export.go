package plan

import (
	"os"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	copyDir "github.com/otiai10/copy"
)

// Export allows save plan to file
func (p *Plan) Export() error {
	if err := os.RemoveAll(p.dir); err != nil {
		return err
	}

	wg := parallel.NewWaitGroup()
	wg.Add(3)

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

func (p *Plan) exportManifest() error {
	if len(p.body.Releases) == 0 {
		return nil
	}

	for k, v := range p.manifests {
		m := filepath.Join(p.dir, Manifest, string(k))

		f, err := helper.CreateFile(m)
		if err != nil {
			return err
		}

		_, err = f.WriteString(v)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Plan) exportGraphMD() error {
	const filename = "graph.md"
	f, err := helper.CreateFile(filepath.Join(p.dir, filename))
	if err != nil {
		return err
	}

	_, err = f.WriteString(p.graphMD)
	if err != nil {
		return err
	}

	return f.Close()
}

func (p *Plan) exportValues() error {
	if len(p.body.Releases) == 0 {
		return nil
	}

	found := false

	for i, rel := range p.body.Releases {
		for j := range p.body.Releases[i].Values {
			found = true
			p.body.Releases[i].Values[j].SetUniq(p.dir, rel.Uniq())
		}
	}

	if !found {
		return nil
	}

	// It doesnt work if workdir is mount.
	//return os.Rename(
	//	filepath.Join(p.tmpDir, Values),
	//	filepath.Join(p.dir, Values),
	//)

	return copyDir.Copy(
		filepath.Join(p.tmpDir, Values),
		filepath.Join(p.dir, Values),
	)
}

// IsExist returns true if planfile exists
func (p *Plan) IsExist() bool {
	return helper.IsExists(p.fullPath)
}

// IsManifestExist returns true if planfile exists
func (p *Plan) IsManifestExist() bool {
	return helper.IsExists(filepath.Join(p.dir, Manifest))
}
