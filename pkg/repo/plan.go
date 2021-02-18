package repo

import (
	log "github.com/sirupsen/logrus"
	"github.com/zhilyaev/helmwave/pkg/helper"
	"github.com/zhilyaev/helmwave/pkg/release"
	"os"
	"strings"
)

func Plan(releases []*release.Config, repositories []*Config) (plan []*Config) {
	all := All(releases)

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
			if _, err := os.Stat(a); !os.IsNotExist(err) {
				found = true
				log.Infof("🗄 %q is local repo", a)
			} else {
				log.Errorf("🗄 %q not found ", a)
			}
		}

	}

	return plan
}

func All(releases []*release.Config) (repos []string) {
	for _, rel := range releases {
		chart := strings.Split(rel.Chart, "/")[0]
		deps, _ := rel.ReposDeps()

		all := append(deps, chart)
		for _, r := range all {
			if !helper.Contains(r, repos) {
				repos = append(repos, r)
			}
		}
	}

	return repos
}
