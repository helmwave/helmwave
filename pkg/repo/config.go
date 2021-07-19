package repo

import (
	"errors"

	"helm.sh/helm/v3/pkg/repo"
)

type Config struct {
	repo.Entry `yaml:",inline"`
	Force      bool
}

var ErrNotFound = errors.New("repository not found")
