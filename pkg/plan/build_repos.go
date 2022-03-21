package plan

import (
	"errors"
	"helm.sh/helm/v3/pkg/registry"
	"os"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) buildRepositories() (out []repo.Config, err error) {
	return buildRepositories(
		buildRepoMapTop(p.body.Releases),
		p.body.Repositories,
	)
}

func buildRepositories(m map[string][]release.Config, in []repo.Config) (out []repo.Config, err error) {
	for rep, releases := range m {
		rm := releaseNames(releases)

		l := log.WithField("repository", rep)
		l.WithField("releases", rm).Debug("🗄 found releases that depend on repository")

		if index, found := repo.IndexOfName(in, rep); found {
			out = append(out, in[index])
			l.Info("🗄 repo has been added to the plan")
		} else {
			l.WithField("releases", rm).Warn("🗄 some releases depend on repository that is not defined")

			return nil, errors.New("🗄 not found " + rep)
		}
	}

	return out, nil
}

// func buildRepoMapDeps(releases []*release.Config) (map[string][]*release.Config, error) {
//	m := make(map[string][]*release.Config)
//	for _, rel := range releases {
//		reps, err := rel.RepoDeps()
//		if err != nil {
//			return nil, err
//		}
//
//		log.WithFields(log.Fields{
//			"release":      rel.Uniq(),
//			"repositories": reps,
//		}).Trace("RepoDeps names")
//
//		for _, rep := range reps {
//			m[rep] = append(m[rep], rel)
//		}
//	}
//
//	return m, nil
//
// }

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

// // allRepos for releases.
// func allRepos(releases []*release.Config) ([]string, error) {
//	var all []string
//	for _, rel := range releases {
//		r, err := rel.RepoDeps()
//		if err != nil {
//			return nil, err
//		}
//
//		all = append(all, r...)
//	}
//
//	return all, nil
// }

// repoIsLocal return true if repo is dir.
func repoIsLocal(repoString string) bool {
	if repoString == "" {
		return true
	}

	stat, err := os.Stat(repoString)
	if (err == nil || !os.IsNotExist(err)) && stat.IsDir() {
		return true
	}

	return false
}
