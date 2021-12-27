package plan

import (
	"fmt"
	"sync"

	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) buildManifest() error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	mu := &sync.Mutex{}

	for _, rel := range p.body.Releases {
		go p.buildReleaseManifest(wg, rel, mu)
	}

	return wg.Wait()
}

func (p *Plan) buildReleaseManifest(wg *parallel.WaitGroup, rel release.Config, mu *sync.Mutex) {
	defer wg.Done()

	if err := rel.ChartDepsUpd(); err != nil {
		log.Warnf("❌ %s cant get dependencies : %v", rel.Uniq(), err)
	}

	rel.DryRun(true)

	r, err := rel.Sync()
	rel.DryRun(false)
	if err != nil || r == nil {
		log.Errorf("❌ %s cant get manifests : %v", rel.Uniq(), err)
		wg.ErrChan() <- err
	}

	hm := ""
	for _, h := range r.Hooks {
		hm += fmt.Sprintf("---\n# Source: %s\n%s\n", h.Path, h.Manifest)
	}

	document := r.Manifest
	if len(r.Hooks) > 0 {
		document += "# ========= HOOKS ========\n" + hm
	}

	log.Trace(document)

	mu.Lock()
	p.manifests[rel.Uniq()] = document
	mu.Unlock()

	log.Infof("✅ %s manifest done", rel.Uniq())
}
