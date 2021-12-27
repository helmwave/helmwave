package repo

import (
	"errors"

	"helm.sh/helm/v3/pkg/repo"
)

type config struct {
	repo.Entry `yaml:",inline"`
	Force      bool
}

func (c *config) Name() string {
	return c.Entry.Name
}

func (c *config) URL() string {
	return c.Entry.URL
}

// ErrNotFound is an error for not declared repository name.
var ErrNotFound = errors.New("repository not found")
