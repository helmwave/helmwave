package plan

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
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
	err = p.buildValues(filepath.Join(os.TempDir(), Dir))
	if err != nil {
		return err
	}

	// Build Repositories
	_, err = buildRepositories(
		buildRepoMapTop(p.body.Releases),
		p.body.Repositories,
	)
	if err != nil {
		return err
	}

	// Sync All
	err = syncRepositories(p.body.Repositories)
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
