package plan

import (
	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"os"
)

// Build plan with yml and tags/matchALL options
func (p *Plan) Build(yml string, tags []string, matchAll bool) error {
	// Create Body
	body, err := NewBody(yml)
	if err != nil {
		return err
	}
	p.body = body

	// Build Releases
	p.body.Releases = buildReleases(tags, p.body.Releases, matchAll)
	if len(p.body.Releases) == 0 {
		return nil
	}

	// Build graph
	p.graphMD = buildGraphMD(p.body.Releases)
	log.Infof("Depends On:\n%s", buildGraphASCII(p.body.Releases))

	// Build Values
	err = p.buildValues(os.TempDir())
	if err != nil {
		return err
	}

	// Build buildRepositories
	reposTop, err := buildRepositories(
		buildRepoMapTop(p.body.Releases),
		p.body.Repositories,
	)

	log.Trace("repo top has been built")
	if err != nil {
		return err
	}
	// Sync Top Repo
	err = syncRepositories(reposTop, helper.Helm)
	if err != nil {
		return err
	}
	log.Trace("repo top has been sync")

	repoMap, err := buildRepoMapDeps(p.body.Releases)
	if err != nil {
		return err
	}
	log.Trace("repo MapAll has been built")

	// Build buildRepositories
	p.body.Repositories, err = buildRepositories(
		repoMap,
		p.body.Repositories,
	)
	log.Trace("repo ALL has been built")
	// Sync All
	err = syncRepositories(reposTop, helper.Helm)
	if err != nil {
		return err
	}

	// Build Manifest
	err = p.buildManifest()
	if err != nil {
		return err
	}

	return nil
}
