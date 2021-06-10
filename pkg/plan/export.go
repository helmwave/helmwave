package plan

import "github.com/helmwave/helmwave/pkg/helper"

// Export allows save plan to file
func (p *Plan) Export() error {
	return helper.Save(p.fullPath, p.body)
}
