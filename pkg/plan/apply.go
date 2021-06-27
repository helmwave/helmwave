package plan

import (
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
	helm "helm.sh/helm/v3/pkg/cli"
)


func (p *Plan) Apply(parallel bool) error {
	if len(p.body.Releases) == 0 {
		return release.ErrEmpty
	}

	log.Info("🛥 Sync releases...")
	//return apply(p.body.Releases, p.dir + PlanManifest, parallel)
	return nil
}



func (p *Plan) SyncRepositories(settings *helm.EnvSettings) (err error) {
	for _, r := range p.body.Repositories {
		err := r.Install(settings)
		if err != nil {
			return err
		}
	}

	return nil
}
