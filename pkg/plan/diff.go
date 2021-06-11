package plan

import "github.com/r3labs/diff/v2"

func (p *Plan) Diff(b *Plan) (diff.Changelog, error) {
	return diff.Diff(p, b, diff.AllowTypeMismatch(true))
}
