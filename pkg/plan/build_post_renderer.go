package plan

import (
	"context"
	"fmt"
	"sync"

	"github.com/helmwave/helmwave/pkg/release"
)

func (p *Plan) buildReleasePostRenderer(ctx context.Context, rel release.Config, mu *sync.Mutex) error {
	templateFuncs := p.templateFuncs(mu)

	if getValuesFunc, ok := templateFuncs["getValues"].(func(string, string) (any, error)); ok {
		templateFuncs["getValues"] = func(args ...string) (any, error) {
			switch len(args) {
			case 1:
				return getValuesFunc(rel.Uniq().String(), args[0])
			case 2:
				return getValuesFunc(args[0], args[1])
			default:
				return nil, fmt.Errorf("getValues requires 1 or 2 arguments")
			}
		}
	}

	return rel.BuildPostRenderer(ctx, p.tmpDir, p.templater, templateFuncs)
}
