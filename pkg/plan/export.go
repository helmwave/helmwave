package plan

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/helper"
	dir "github.com/otiai10/copy"
)

type Exporter interface {
	Export(p *Plan) error
}

type ExporterS3 struct {
}

func (e *ExporterS3) Export(p *Plan) error {
	return nil
}

type ExporterLocal struct {
}

func (e *ExporterLocal) Export(p *Plan) error {
	return nil
}

var Exporters = map[string]Exporter{
	"s3://":   &ExporterS3{},
	"file://": &ExporterLocal{},
}

// Export allows save plan to file.
func (p *Plan) Export() error {
	exporter, found := Exporters[p.url.Scheme]
	if !found {
		return fmt.Errorf("plan export to '%s' is not supported", p.url.Scheme)
	}

	return exporter.Export(p)
}

func (p *Plan) exportManifest() error {
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

func (p *Plan) exportGraphMD() (err error) {
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

	const filename = "graph.md"
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

func (p *Plan) exportValues() error {
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

	// It doesnt work if workdir has been mounted.
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

// IsExist returns true if planfile exists.
func (p *Plan) IsExist() bool {
	return helper.IsExists(p.File())
}

// IsManifestExist returns true if planfile exists.
func (p *Plan) IsManifestExist() bool {
	return helper.IsExists(filepath.Join(p.Dir(), Manifest))
}
