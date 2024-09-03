package plan

import "github.com/helmwave/helmwave/pkg/release/uniqname"

func (p *Plan) Manifests() map[uniqname.UniqName]string {
	return p.manifests
}
