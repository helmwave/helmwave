package repo

import (
	"github.com/zhilyaev/helmwave/pkg/helper"
	"github.com/zhilyaev/helmwave/pkg/release"
	"strings"
)

func Plan(releases []release.Config, repositories []Config) (plan []Config) {
	for i := len(repositories) - 1; i >= 0; i-- {
		for j := len(releases) - 1; j >= 0; j-- {
			// bitnami/redis -> bitnami
			name := strings.Split(releases[j].Chart, "/")[0]
			deps, _ := releases[j].ReposDeps()

			if (name == repositories[i].Name || helper.Contains(repositories[i].Name, deps)) && !repositories[i].In(plan) {
				plan = append(plan, repositories[i])
				releases = append(releases[:j], releases[j+1:]...)
				break
			}
		}
	}
	return plan
}
