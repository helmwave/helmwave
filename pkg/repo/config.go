package repo

import (
	"errors"
	"helm.sh/helm/v3/pkg/repo"
)

type Config struct {
	repo.Entry `yaml:",inline"`

	// TODO: Support Flag
	Force bool
}

var ErrNotFound = errors.New("repository not found")
