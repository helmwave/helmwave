package plan

import (
	"fmt"

	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) buildManifest() error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for _, rel := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel *release.Config) {
			defer wg.Done()

			err := rel.ChartDepsUpd()
			if err != nil {
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

			m := rel.Uniq() + ".yml"
			p.manifests[m] = document

			log.Infof("✅ %s manifest done", rel.Uniq())
		}(wg, rel)
	}

	return wg.Wait()
}
