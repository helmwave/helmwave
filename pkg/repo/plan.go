package repo

import (
	"os"
	"strings"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

// Plan generates repo config out of planned releases and available repositories
func Plan(releases []*release.Config, repositories []*Config) (plan []*Config) {
	all := getRepositories(releases)

	for _, a := range all {
		found := false
		for _, b := range repositories {
			if a == b.Name {
				found = true
				if !b.InByName(plan) {
					plan = append(plan, b)
					log.Infof("🗄 %q has been added to the plan", a)
				}
			}
		}

		if !found {
			log.Errorf("🗄 %q not found ", a)
		}
	}

	return plan
}

// Get repositories for releases
func getRepositories(releases []*release.Config) (repos []string) {
	for _, rel := range releases {
		repo := strings.Split(rel.Chart, "/")[0]
		deps, _ := rel.ReposDeps()

		all := deps
		if repoIsLocal(repo) {
			log.Infof("🗄 %q is local repo", repo)
		} else {
			all = append(all, repo)
		}

		for _, r := range all {
			if !helper.Contains(r, repos) {
				repos = append(repos, r)
			}
		}
	}

	return repos
}

func repoIsLocal(repo string) bool {
	if repo == "" {
		return true
	}

	stat, err := os.Stat(repo)
	if (err == nil || !os.IsNotExist(err)) && stat.IsDir() {
		return true
	}

	return false
}
