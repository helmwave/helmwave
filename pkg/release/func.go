package release

import (
	"github.com/zhilyaev/helmwave/pkg/helper"
	"github.com/zhilyaev/helmwave/pkg/template"
	"helm.sh/helm/v3/pkg/chart/loader"
	"strings"
)

func (rel *Config) In(a []Config) bool {
	for _, r := range a {
		if rel == &r {
			return true
		}
	}
	return false
}

func (rel *Config) RenderValues(dir string) error {
	rel.PlanValues()

	for i, v := range rel.Values {

		s := v + "." + rel.Name + "@" + rel.Options.Namespace + ".plan"

		p := dir + s
		err := template.Tpl2yml(v, p, struct{ Release *Config }{rel})
		if err != nil {
			return err
		}

		rel.Values[i] = p
	}

	return nil
}

func (rel *Config) ReposDeps() (repos []string, err error) {
	chart, err := loader.Load(rel.Chart)
	if err != nil {
		return nil, err
	}

	deps := chart.Metadata.Dependencies

	for _, d := range deps {
		if strings.HasPrefix(d.Repository, "@") {
			d.Repository = helper.TrimFirstRune(d.Repository)
		}
		repos = append(repos, d.Repository)
	}

	return repos, nil
}
