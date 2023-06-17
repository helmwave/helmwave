package plan

import (
	"context"

	"github.com/helmwave/helmwave/pkg/hooks"
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
func (p *Plan) Build(ctx context.Context, o BuildOptions) error { //nolint:funlen,cyclop
	p.templater = o.Templater

	// Create Body
	body, err := NewBody(ctx, o.Yml)
	if err != nil {
		return err
	}
	p.body = body

	// Run pre-build hooks
	if len(p.body.Hooks.PreBuild) != 0 {
		log.Info("🩼 Running pre-build hooks...")
		hooks.Run(p.body.Hooks.PreBuild)
	}

	// Build Releases
	log.Info("🔨 Building releases...")
	p.body.Releases, err = buildReleases(o.Tags, p.body.Releases, o.MatchAll)
	if err != nil {
		return err
	}
	if len(p.body.Releases) == 0 {
		return nil
	}

	// Build graphs
	if o.GraphWidth != 1 {
		log.Info("🔨 Building graphs...")
		p.graphMD = buildGraphMD(p.body.Releases)
		log.Infof("show graph:\n%s", p.BuildGraphASCII(o.GraphWidth))
	}

	// Build Values
	log.Info("🔨 Building values...")
	err = p.buildValues()
	if err != nil {
		return err
	}

	// Build Repositories
	log.Info("🔨 Building repositories...")
	_, err = p.buildRepositories()
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
	_, err = p.buildRegistries()
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
	err = p.buildCharts()
	if err != nil {
		return err
	}

	// Build Manifest
	log.Info("🔨 Building manifests...")
	err = p.buildManifest(ctx)
	if err != nil {
		return err
	}

	// Run post-build hooks
	if len(p.body.Hooks.PostBuild) != 0 {
		log.Info("🩼 Running post-build hooks...")
		hooks.Run(p.body.Hooks.PostBuild)
	}

	return nil
}
