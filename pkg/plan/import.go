package plan

import (
	"fmt"
)

// Import parses directory with plan files and imports them into structure.
func (p *Plan) Import() error {
	b, found := Backends[p.url.Scheme]
	if !found {
		return fmt.Errorf("plan import to %q is not supported", p.url.Scheme)
	}

	return b.Import(p)
}
