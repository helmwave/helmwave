// This package exports some fields for tests
package repo

import "helm.sh/helm/v3/pkg/repo"

func NewConfig() *config {
	return &config{
		Entry: repo.Entry{
			Name: "bla",
			URL:  "https://bla",
		},
	}
}
