package plan

import (
	"github.com/helmwave/helmwave/pkg/helper"
	"os"
)

// Export allows save plan to file
func (p *Plan) Export() error {
	if err := os.RemoveAll(p.dir); err != nil {
		return err
	}

	if err := p.exportManifest(); err != nil {
		return err
	}

	// Todo: Rename tmpDir is faster than rebuild values
	if err := p.buildValues(p.dir); err != nil {
		return err
	}

	if err := p.exportGraphMD(); err != nil {
		return err
	}

	return helper.SaveInterface(p.fullPath, p.body)
}

func (p *Plan) exportManifest() error {
	for k, v := range p.manifests {
		m := p.dir + Manifest + string(k)

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
	f, err := helper.CreateFile(p.dir + "graph.md")
	if err != nil {
		return err
	}

	_, err = f.WriteString(p.graphMD)
	if err != nil {
		return err
	}

	return f.Close()
}

// IsExist returns true if planfile exists
func (p *Plan) IsExist() bool {
	return helper.IsExists(p.fullPath)
}

//IsManifestExist returns true if planfile exists
func (p *Plan) IsManifestExist() bool {
	return helper.IsExists(p.dir + Manifest)
}
