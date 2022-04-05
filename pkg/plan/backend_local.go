package plan

import (
	"errors"
	"fmt"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/version"
	dir "github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type BackendLocal struct{}

func (e *BackendLocal) Import(p *Plan) error {
	body, err := NewBody(p.fsys, p.File())
	if err != nil {
		return err
	}

	err = e.importManifest(p)
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

func (e *BackendLocal) importManifest(p *Plan) error {
	d := filepath.Join(p.Dir(), Manifest)
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

func (e *BackendLocal) Export(p *Plan) error {
	if err := os.RemoveAll(p.Dir()); err != nil {
		return fmt.Errorf("failed to clean plan directory %s: %w", p.Dir(), err)
	}

	wg := parallel.NewWaitGroup()
	wg.Add(3)

	go func() {
		defer wg.Done()
		if err := e.exportManifests(p); err != nil {
			wg.ErrChan() <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := e.exportValues(p); err != nil {
			wg.ErrChan() <- err
		}

		// Save Planfile after values
		if err := helper.SaveInterface(p.File(), p.body); err != nil {
			wg.ErrChan() <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := e.exportGraphMD(p); err != nil {
			wg.ErrChan() <- err
		}
	}()

	return wg.Wait()
}

func (e *BackendLocal) exportManifests(p *Plan) error {
	if len(p.body.Releases) == 0 {
		return nil
	}

	for k, v := range p.manifests {
		m := filepath.Join(p.Dir(), Manifest, string(k)+".yml")

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

// IsExist returns true if planfile exists.
func (e *BackendLocal) IsExist(p *Plan) bool {
	return helper.IsExists(p.File())
}

// IsManifestExist returns true if planfile exists.
func (e *BackendLocal) IsManifestExist(p *Plan) bool {
	return helper.IsExists(filepath.Join(p.Dir(), Manifest))
}

func (e *BackendLocal) exportValues(p *Plan) error {
	if len(p.body.Releases) == 0 {
		return nil
	}

	found := false

	for i, rel := range p.body.Releases {
		for j := range p.body.Releases[i].Values() {
			found = true
			p.body.Releases[i].Values()[j].SetUniq(p.Dir(), rel.Uniq())
		}
	}

	if !found {
		return nil
	}

	// It doesn't work if workdir has been mounted.
	err := os.Rename(
		filepath.Join(p.tmpDir, Values),
		filepath.Join(p.Dir(), Values),
	)
	if err != nil {
		err = dir.Copy(
			filepath.Join(p.tmpDir, Values),
			filepath.Join(p.Dir(), Values),
		)
		if err != nil {
			return fmt.Errorf("failed to copy values from %s to %s: %w", p.tmpDir, p.Dir(), err)
		}

		return nil
	}

	return nil
}

func (e *BackendLocal) exportGraphMD(p *Plan) (err error) {
	if len(p.body.Releases) == 0 {
		return nil
	}

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

	//f, err := helper.CreateFile(filepath.Join(p.URL.Path, filename))

	f, err := p.fsys.Open(filepath.Join(p.Dir(), filename))

	if err != nil {
		return err
	}

	//_, err = f.WriteString(p.graphMD)
	//if err != nil {
	//	return fmt.Errorf("failed to write graph file %s: %w", f.Name(), err)
	//}

	if err = f.Close(); err != nil {
		return fmt.Errorf("failed to close graph file %s: %w", filename, err)
	}

	return nil
}
