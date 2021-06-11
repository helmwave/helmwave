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
