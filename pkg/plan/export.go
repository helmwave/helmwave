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
	return helper.Save(p.fullPath, p.body)
}

// IsExist returns true if planfile exists
func (p *Plan) IsExist() bool {
	return  helper.IsExists(p.fullPath)
}

//IsManifestExist returns true if planfile exists
func (p *Plan) IsManifestExist() bool {
	return  helper.IsExists(p.dir + PlanManifest)
}
