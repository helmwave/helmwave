package plan

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type BuildOptions struct { //nolint:govet
	Tags       []string
	Yml        string
	Templater  string
	MatchAll   bool
	GraphWidth int
}

func (o *BuildOptions) Body(ctx context.Context) (body *planBody, err error) {
	// Create Body
	body, err = NewBody(ctx, o.Yml, false)
	if err != nil {
		return
	}
	return body, nil
}

// Build plan with yml and tags/matchALL options.
//
//nolint:cyclop,gocognit // TODO: reduce cyclomatic complexity
func (p *Plan) Build(ctx context.Context, o BuildOptions) (err error) { //nolint:funlen
	p.templater = o.Templater
	p.body, err = o.Body(ctx)
	if err != nil {
		return
	}
	return p.build(ctx, o)

}

//nolint:cyclop,gocognit
func (p *Plan) BuildWithBody(ctx context.Context, o BuildOptions, body *planBody) (err error) { //nolint:funlen
	p.body = body
	return p.build(ctx, o)

}

func (p *Plan) build(ctx context.Context, o BuildOptions) (err error) {
	// Run hooks
	err = p.body.Lifecycle.RunPreBuild(ctx)
	if err != nil {
		return
	}

	defer func() {
		lifecycleErr := p.body.Lifecycle.RunPostBuild(ctx)
		if lifecycleErr != nil {
			log.Errorf("got an error from postbuild hooks: %v", lifecycleErr)
			if err == nil {
				err = lifecycleErr
			}
		}
	}()

	// Build Releases
	log.Info("🔨 Building releases...")
	p.body.Releases, err = p.buildReleases(o.Tags, o.MatchAll)
	if err != nil {
		return
	}

	// Build Values
	log.Info("🔨 Building values...")
	err = p.buildValues(ctx)
	if err != nil {
		return
	}

	// Build Repositories
	log.Info("🔨 Building repositories...")
	p.body.Repositories, err = p.buildRepositories()
	if err != nil {
		return
	}

	// Sync Repositories
	err = SyncRepositories(ctx, p.body.Repositories)
	if err != nil {
		return
	}

	// Build Registries
	log.Info("🔨 Building registries...")
	p.body.Registries, err = p.buildRegistries()
	if err != nil {
		return
	}
	// Sync Registries
	err = p.syncRegistries(ctx)
	if err != nil {
		return
	}

	// to build charts, we need repositories and registries first
	log.Info("🔨 Building charts...")
	err = p.buildCharts()
	if err != nil {
		return
	}

	// Validating plan after it was changed
	err = p.body.Validate()
	if err != nil {
		return
	}

	// Build Manifest
	log.Info("🔨 Building manifests...")
	err = p.buildManifest(ctx)
	if err != nil {
		return
	}

	// Build graphs
	if o.GraphWidth != 1 {
		log.Info("🔨 Building graphs...")
		p.graphMD = buildGraphMD(p.body.Releases)
		log.Infof("show graph:\n%s", p.BuildGraphASCII(o.GraphWidth))
	}

	return
}
