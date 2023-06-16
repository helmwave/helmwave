package plan

import (
	"errors"

	regi "github.com/helmwave/helmwave/pkg/registry"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/registry"
)

func buildRegistriesMapTop(releases []release.Config) map[string][]release.Config {
	m := make(map[string][]release.Config)
	for _, rel := range releases {
		if registry.IsOCI(rel.Chart().Name) {
			m[rel.Repo()] = append(m[rel.Repo()], rel)
			rel.Logger().Debugln("🗄 This chart will download via OCI")
		}
	}

	return m
}

func buildRegistries(m map[string][]release.Config, in []regi.Config) (out []regi.Config, err error) {
	for reg, releases := range m {
		rm := releaseNames(releases)

		l := log.WithField("registry", reg)
		l.WithField("releases", rm).Debug("🗄 found releases that depend on registries")

		if index, found := regi.IndexOfHost(in, reg); found {
			out = append(out, in[index])
			l.Info("🗄 registry has been added to the plan")
		} else {
			l.WithField("releases", rm).Warn("🗄 some releases depend on repository that is not defined")

			return nil, errors.New("🗄 not found " + reg)
		}
	}

	return out, nil
}

func (p *Plan) buildRegistries() (out []regi.Config, err error) {
	return buildRegistries(
		buildRegistriesMapTop(p.body.Releases),
		p.body.Registries,
	)
}
