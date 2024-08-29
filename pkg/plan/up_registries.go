package plan

import (
	"context"

	"github.com/helmwave/helmwave/pkg/parallel"
	regi "github.com/helmwave/helmwave/pkg/registry"
)

func (p *Plan) syncRegistries(ctx context.Context) (err error) {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Registries))

	for i := range p.body.Registries {
		go func(wg *parallel.WaitGroup, reg regi.Config) {
			defer wg.Done()
			err := reg.Install()
			if err != nil {
				wg.ErrChan() <- err
			}
		}(wg, p.body.Registries[i])
	}

	if err := wg.WaitWithContext(ctx); err != nil {
		return err
	}

	return err
}
