package plan

import (
	"fmt"
)

// Export save plan via backend
func (p *Plan) Export() error {
	b, found := Backends[p.url.Scheme]
	if !found {
		return fmt.Errorf("plan export to %q is not supported", p.url.Scheme)
	}

	return b.Export(p)
}
