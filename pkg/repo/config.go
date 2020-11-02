package repo

import (
	"helm.sh/helm/v3/pkg/repo"
)

type Config struct {
	repo.Entry `yaml:",inline"`

	// TODO: Support Flag
	Force bool
}
