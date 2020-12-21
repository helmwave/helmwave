package repo

import (
	"errors"
	"github.com/zhilyaev/helmwave/pkg/release"
	"strings"
)

func Plan(releases []release.Config, repositories []Config) (plan []Config, err error) {
	for j := len(releases) - 1; j >= 0; j-- {
		// bitnami/redis -> bitnami
		chart := strings.Split(releases[j].Chart, "/")[0]
		deps, _ := releases[j].ReposDeps()
		reps := append(deps, chart)

		for i := len(reps) - 1; i >= 0; i-- {
			for k := len(repositories) - 1; k >= 0; k-- {
				if reps[i] == repositories[k].Name {
					if !repositories[i].In(plan) {
						plan = append(plan, repositories[i])
						repositories = append(repositories[:i], repositories[i:]...)
					}
					continue
				}
			}
			err = errors.New(reps[i] + " not found in the repositories")
			return plan, err
		}
	}

	return plan, nil
}
