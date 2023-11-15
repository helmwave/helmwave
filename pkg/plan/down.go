package plan

import (
	"context"

	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/dependency"
	log "github.com/sirupsen/logrus"
)

// Down destroys all releases that exist in a plan.
func (p *Plan) Down(ctx context.Context) error {
	dependenciesGraph, err := p.body.generateDependencyGraph()
	if err != nil {
		return err
	}

	dependenciesGraph, err = dependenciesGraph.Reverse()
	if err != nil {
		return err
	}

	// Run hooks
	err = p.body.Lifecycle.RunPreDown(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err := p.body.Lifecycle.RunPostDown(ctx)
		if err != nil {
			log.Errorf("got an error from postdown hooks: %v", err)
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
			node.SetSucceeded()
			if err != nil {
				log.Errorf("❌ %s: %v", rel.Uniq(), err)
				wg.ErrChan() <- err
			} else {
				log.Infof("✅ %s uninstalled!", rel.Uniq())
			}
		}(ctx, wg, node)
	}

	err = wg.Wait()
	if err != nil {
		return err
	}

	return nil
}
