package action

import (
	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
)

type Up struct {
	build *Build
	dog   *kubedog.Config

	autoBuild      bool
	kubedogEnabled bool
}

func (i *Up) Run() error {
	if i.autoBuild {
		if err := i.build.Run(); err != nil {
			return err
		}
	}

	p := plan.New(i.build.plandir)
	if err := p.Import(); err != nil {
		return err
	}

	p.PrettyPlan()

	if i.kubedogEnabled {
		log.Warn("🐶 kubedog is enable")
		return p.ApplyWithKubedog(i.dog)
	}

	return p.Apply()
}
