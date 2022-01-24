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

	l := log.WithField("release", rel.Uniq())

	if err := rel.ChartDepsUpd(); err != nil {
		l.Warnf("❌ can't get dependencies : %v", err)
	}

	rel.DryRun(true)

	r, err := rel.Sync()
	rel.DryRun(false)
	if err != nil || r == nil {
		l.Errorf("❌ can't get manifests: %v", err)
		wg.ErrChan() <- err

		return
	}

	hm := ""
	for _, h := range r.Hooks {
		hm += fmt.Sprintf("---\n# Source: %s\n%s\n", h.Path, h.Manifest)
	}

	document := r.Manifest
	if len(r.Hooks) > 0 {
		document += "# ========= HOOKS ========\n" + hm
	}

	l.Trace(document)

	mu.Lock()
	p.manifests[rel.Uniq()] = document
	mu.Unlock()

	l.Info("✅ manifest done")
}
