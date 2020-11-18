package release

import (
	"fmt"
	"github.com/zhilyaev/helmwave/pkg/helper"
	"github.com/zhilyaev/helmwave/pkg/template"
	"helm.sh/helm/v3/pkg/chart/loader"
	"os"
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

func (rel *Config) PlanValues() {

	for i := len(rel.Values) - 1; i >= 0; i-- {
		if _, err := os.Stat(rel.Values[i]); err != nil {
			if os.IsNotExist(err) {
				rel.Values = append(rel.Values[:i], rel.Values[i+1:]...)
			}
		}
	}

}

func (rel *Config) RenderValues(debug bool) {
	rel.PlanValues()

	for i, v := range rel.Values {
		p := v + ".plan"
		err := template.Tpl2yml(v, p, struct{ Release *Config }{rel}, debug)
		if err != nil {
			fmt.Println(err)
		}

		rel.Values[i] = p
	}

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
