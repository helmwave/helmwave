package plan

import (
	"context"
	"io/fs"

	"github.com/helmwave/go-fsimpl"
	log "github.com/sirupsen/logrus"
)

type BuildOptions struct { //nolint:govet
	Tags       []string
	Yml        string
	Templater  string
	MatchAll   bool
	GraphWidth int
}

// Build plan with yml and tags/matchALL options.
//
//nolint:cyclop // TODO: reduce cyclomatic complexity
func (p *Plan) Build(ctx context.Context, srcFS fs.StatFS, destFS fsimpl.WriteableFS, o BuildOptions) error { //nolint:funlen
	p.templater = o.Templater

	// Create Body
	body, err := NewBody(ctx, srcFS, o.Yml, false)
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

	// Build Values
	log.Info("🔨 Building values...")
	err = p.buildValues(srcFS, destFS)
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

	// to build charts, we need repositories and registries first
	log.Info("🔨 Building charts...")
	err = p.buildCharts(srcFS, destFS)
	if err != nil {
		return err
	}

	// Validating plan after it was changed
	err = p.body.Validate()
	if err != nil {
		return err
	}

	// Build Manifest
	log.Info("🔨 Building manifests...")
	err = p.buildManifest(ctx, srcFS)
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
