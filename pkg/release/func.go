package release

import (
	"fmt"
	"github.com/helmwave/helmwave/pkg/template"
	"helm.sh/helm/v3/pkg/chart/loader"
	"path"
	"strings"
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
		m, err := v.ManifestPath()
		if err != nil {
			return err
		}
		s := fmt.Sprintf("%s.%s.plan", m, rel.UniqName())
		p := path.Join(dir, s)

		err = template.Tpl2yml(v.GetPath(), p, struct{ Release *Config }{rel})
		v.UnlinkProcessed()
		if err != nil {
			return err
		}

		rel.Values[i].SetProcessedPath(p)
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
		d.Repository = strings.TrimPrefix(d.Repository, "@")
		repos = append(repos, d.Repository)
	}

	return repos, nil
}
