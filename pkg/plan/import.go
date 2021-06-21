package plan

import "github.com/helmwave/helmwave/pkg/version"

func (p *Plan) Import() error {
	body, err := NewBody(p.fullPath)
	if err != nil {
		return err
	}

	p.body = body
	version.Check(p.body.Version, version.Version)

	return nil
}
