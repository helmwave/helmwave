package plan

import (
	"context"
	"io/fs"

	log "github.com/sirupsen/logrus"
)

type BuildOptions struct { //nolint:govet
	Tags       []string
	Yml        fs.FS
	Templater  string
	MatchAll   bool
	GraphWidth int
}

// Build plan with yml and tags/matchALL options.
func (p *Plan) Build(ctx context.Context, o BuildOptions) error {
	p.templater = o.Templater

	// Create Body
	body, err := NewBody(ctx, o.Yml, false)
	if err != nil {
		return err
	}
	p.body = body

	// Run hooks
	err = p.body.Lifecycle.RunPreBuild(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err := p.body.Lifecycle.RunPostBuild(ctx)
		if err != nil {
			log.Errorf("got an error from postbuild hooks: %v", err)
		}
	}()

	// Build Releases
	log.Info("🔨 Building releases...")
	p.body.Releases, err = p.buildReleases(o.Tags, o.MatchAll)
	if err != nil {
		return err
	}

	// Build Repositories
	log.Info("🔨 Building repositories...")
	p.body.Repositories, err = p.buildRepositories()
	if err != nil {
		return err
	}

	// Sync Repositories
	err = SyncRepositories(ctx, p.body.Repositories)
	if err != nil {
		return err
	}

	// Build Registries
	log.Info("🔨 Building registries...")
	p.body.Registries, err = p.buildRegistries()
	if err != nil {
		return err
	}
	// Sync Registries
	err = p.syncRegistries(ctx)
	if err != nil {
		return err
	}

	// Validating plan after it was changed
	err = p.body.Validate()
	if err != nil {
		return err
	}

	// Build graphs
	if o.GraphWidth != 1 {
		log.Info("🔨 Building graphs...")
		p.graphMD = buildGraphMD(p.body.Releases)
		log.Infof("show graph:\n%s", p.BuildGraphASCII(o.GraphWidth))
	}

	return nil
}
