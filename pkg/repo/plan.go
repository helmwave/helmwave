package repo

import (
	"github.com/zhilyaev/helmwave/pkg/helper"
	"github.com/zhilyaev/helmwave/pkg/release"
	"strings"
)

func Plan(releases []release.Config, repositories []Config) []Config {
	var plan []Config

	for _, rep := range repositories {
		for _, rel := range releases {
			// bitnami/redis -> bitnami
			name := strings.Split(rel.Chart, "/")[0]
			deps, _ := rel.ReposDeps()

			if (name == rep.Name || helper.Contains(rep.Name, deps)) && !rep.In(plan) {
				plan = append(plan, rep)
				//release.RemoveIndex(releases, i) // Optimise deleter
				break
			}

		}
	}

	return plan
}
