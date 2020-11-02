package repo

import (
	"github.com/zhilyaev/helmwave/pkg/release"
	"strings"
)

func Plan(releases []release.Config, repositories []Config) []Config {
	var plan []Config

	for _, rep := range repositories {
		for _, rel := range releases {
			// bitnami/redis -> bitnami
			name := strings.Split(rel.Chart, "/")[0]
			if name == rep.Name && !rep.In(plan) {
				plan = append(plan, rep)
				//release.RemoveIndex(releases, i)
				break
			}
		}
	}

	return plan
}
