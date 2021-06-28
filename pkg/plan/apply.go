package plan

import (
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
	helm "helm.sh/helm/v3/pkg/cli"
)

func (p *Plan) Apply(parallel bool) (err error) {
	if len(p.body.Releases) == 0 {
		return release.ErrEmpty
	}

	log.Info("🗄 Sync repositories...")
	err = p.syncRepositories(helm.New())
	if err != nil {
		return err
	}

	log.Info("🛥 Sync releases...")
	err = p.syncReleases()
	if err != nil {
		return err
	}

	return nil
}

func (p *Plan) syncRepositories(settings *helm.EnvSettings) (err error) {
	for _, r := range p.body.Repositories {
		err := r.Install(settings)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Plan) syncReleases() (err error) {
	for _, r := range p.body.Releases {
		_, err = r.Sync()
		if err != nil {
			return err
		}
	}

	return nil
}
