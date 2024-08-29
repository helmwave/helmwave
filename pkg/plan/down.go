package plan

import (
	"context"

	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/dependency"
	log "github.com/sirupsen/logrus"
)

// Down destroys all releases that exist in a plan.
func (p *Plan) Down(ctx context.Context) (err error) {
	dependenciesGraph, err := p.Graph().Reverse()
	if err != nil {
		return err
	}

	// Run hooks
	err = p.body.Lifecycle.RunPreDown(ctx)
	if err != nil {
		return err
	}

	defer func() {
		lifecycleErr := p.body.Lifecycle.RunPostDown(ctx)
		if lifecycleErr != nil {
			log.Errorf("got an error from postdown hooks: %v", lifecycleErr)
			if err == nil {
				err = lifecycleErr
			}
		}
	}()

	nodesChan := dependenciesGraph.Run()

	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for node := range nodesChan {
		go func(ctx context.Context, wg *parallel.WaitGroup, node *dependency.Node[release.Config]) {
			defer wg.Done()
			rel := node.Data
			_, err := rel.Uninstall(ctx)
			if err != nil {
				log.Errorf("❌ %s: %v", rel.Uniq(), err)
				wg.ErrChan() <- err
				node.SetFailed()
			} else {
				log.Infof("✅ %s uninstalled!", rel.Uniq())
				node.SetSucceeded()
			}
		}(ctx, wg, node)
	}

	err = wg.Wait()

	return err
}
