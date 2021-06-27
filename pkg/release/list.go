package release

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"regexp"
)


func (rel *Config) List() (*release.Release, error) {
	cfg, err := rel.cfg()
	if err != nil {
		return nil, err
	}

	client := action.NewList(cfg)
	client.Filter = fmt.Sprintf("^%s$", regexp.QuoteMeta(rel.ReleaseName))

	result, err := client.Run()
	if err != nil {
		return nil, err
	}


	switch len(result) {
	case 0:
		return nil, ErrNotFound
	case 1:
		return result[0], nil
	default:
		return nil, ErrFoundMultiple
	}
}