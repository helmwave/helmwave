package plan

import (
	"io/fs"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/registry"
)

func (p *Plan) buildRepositories(baseFS fs.FS) (out []repo.Config, err error) {
	return buildRepositories(
		baseFS,
		buildRepoMapTop(p.body.Releases),
		p.body.Repositories,
	)
}

func buildRepositories(baseFS fs.FS, m map[string][]release.Config, in []repo.Config) (out []repo.Config, err error) {
	for rep, releases := range m {
		rm := releaseNames(releases)

		l := log.WithField("repository", rep)
		l.WithField("releases", rm).Debug("🗄 found releases that depend on repository")

		if repoIsLocal(baseFS, rep) {
			l.Info("🗄 it is local repo")
		} else if index, found := repo.IndexOfName(in, rep); found {
			out = append(out, in[index])
			l.Info("🗄 repo has been added to the plan")
		} else {
			l.WithField("releases", rm).Warn("🗄 some releases depend on repository that is not defined")

			return nil, repo.NewNotFoundError(rep)
		}
	}

	return out, nil
}

func buildRepoMapTop(releases []release.Config) map[string][]release.Config {
	m := make(map[string][]release.Config)
	for _, rel := range releases {
		// Added to map if is not OCI
		if !registry.IsOCI(rel.Chart().Name) {
			m[rel.Repo()] = append(m[rel.Repo()], rel)
		}
	}

	return m
}

// repoIsLocal returns true if repo is a dir.
func repoIsLocal(baseFS fs.FS, repoString string) bool {
	if repoString == "" {
		return true
	}

	return helper.IsExists(baseFS, repoString) && helper.IsDir(baseFS, repoString)
}
