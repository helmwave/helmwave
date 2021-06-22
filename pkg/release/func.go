package release

import (
	"os"
	"strings"

	"github.com/helmwave/helmwave/pkg/template"
	"helm.sh/helm/v3/pkg/chart/loader"
)

func (rel *Config) UniqName() string {
	return rel.Name + "@" + rel.Options.Namespace
}

func (rel *Config) In(a []*Config) bool {
	for _, r := range a {
		if rel == r {
			return true
		}
	}
	return false
}

func (rel *Config) RenderValues(dir string) error {
	rel.filterValuesFiles()

	for i, v := range rel.Values {

		s := v + "." + rel.UniqName() + ".plan"

		p := dir + s
		err := template.Tpl2yml(v, p, struct{ Release *Config }{rel})
		if err != nil {
			return err
		}

		rel.Values[i] = p
	}

	return nil
}

// filterValuesFiles filters non-existent values files.
func (rel *Config) filterValuesFiles() {
	for i := len(rel.Values) - 1; i >= 0; i-- {
		stat, err := os.Stat(rel.Values[i])
		if os.IsNotExist(err) || stat.IsDir() {
			rel.Values = append(rel.Values[:i], rel.Values[i+1:]...)
		}
	}
}

func (rel *Config) ReposDeps() (repos []string, err error) {
	chart, err := loader.Load(rel.Chart)
	if err != nil {
		return nil, err
	}

	deps := chart.Metadata.Dependencies

	for _, d := range deps {
		d.Repository = strings.TrimPrefix(d.Repository, "@")
		repos = append(repos, d.Repository)
	}

	return repos, nil
}
