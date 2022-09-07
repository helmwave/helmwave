package release

import (
	"strings"

	"helm.sh/helm/v3/pkg/registry"
)

// // RepoDeps returns repository for release
// func (rel *Config) RepoDeps() (repos []string, err error) {
//	chart, err := rel.GetChart()
//	if err != nil {
//		log.Warn("Failed GetChart ", rel.Uniq())
//		return nil, err
//	}
//
//	repos = append(repos, rel.Repo())
//
//	for _, d := range chart.Metadata.Dependencies {
//		if d.Enabled {
//			d.Repository = strings.TrimPrefix(d.Repository, "@")
//			repos = append(repos, d.Repository)
//		} else if helper.IsURL(d.Repository) {
//			return nil, errors.New("url is unsupported: " + d.Repository)
//		}
//	}
//
//	return repos, nil
// }

func (rel *Release) Repo() string {
	s := rel.Chart().Name
	if registry.IsOCI(s) {
		s = strings.TrimPrefix(s, registry.OCIScheme+"://")
	}

	return strings.Split(s, "/")[0]
}
