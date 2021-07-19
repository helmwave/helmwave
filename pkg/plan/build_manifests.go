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

			rel.DryRun(true)
			r, err := rel.Sync()
			rel.DryRun(false)
			if err != nil || r == nil {
				log.Error("i cant generate manifest for ", rel.Uniq())
				log.Fatal(err)
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

			// log.Debug(rel.Uniq(), "`s manifest was successfully built ")
		}(wg, rel)
	}

	return wg.Wait()
}
