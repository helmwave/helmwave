// This package exports some fields for tests
package repo

import "helm.sh/helm/v4/pkg/repo/v1"

func NewConfig() *config {
	return &config{
		Entry: repo.Entry{
			Name: "bla",
			URL:  "https://bla",
		},
	}
}
