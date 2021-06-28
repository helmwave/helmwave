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

	for k, v := range p.manifests {
		m := p.dir + Manifest + k

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

	return helper.SaveInterface(p.fullPath, p.body)
}

// IsExist returns true if planfile exists
func (p *Plan) IsExist() bool {
	return helper.IsExists(p.fullPath)
}

//IsManifestExist returns true if planfile exists
func (p *Plan) IsManifestExist() bool {
	return helper.IsExists(p.dir + Manifest)
}
