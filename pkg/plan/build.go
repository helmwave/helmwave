package plan

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type BuildOptions struct {
	Yml                string
	Templater          string
	Tags               []string
	GraphWidth         int
	MatchAll           bool
	EnableDependencies bool
}

// Build plan with yml and tags/matchALL options.
func (p *Plan) Build(ctx context.Context, o BuildOptions) (err error) {
	p.templater = o.Templater

	// Create Body
	var body *planBody
	body, err = NewBody(ctx, o.Yml, false)
	if err != nil {
		return
	}
	p.body = body

	// Run Pre hooks
	err = p.body.Lifecycle.RunPreBuild(ctx)
	if err != nil {
		return
	}

	err = p.build(ctx, o)
	if err != nil {
		return
	}

	// Run Post hooks
	err = p.body.Lifecycle.RunPostBuild(ctx)
	if err != nil {
		return
	}

	return
}

func (p *Plan) build(ctx context.Context, o BuildOptions) error {
	var err error

	p.body.Releases, err = p.buildReleases(ctx, o)
	if err != nil {
		return err
	}

	err = p.buildValues(ctx)
	if err != nil {
		return err
	}

	p.body.Repositories, err = p.buildRepositories()
	if err != nil {
		return err
	}

	err = SyncRepositories(ctx, p.body.Repositories)
	if err != nil {
		return err
	}

	p.body.Registries, err = p.buildRegistries()
	if err != nil {
		return err
	}

	err = p.syncRegistries(ctx)
	if err != nil {
		return err
	}

	err = p.buildCharts()
	if err != nil {
		return err
	}

	err = p.body.Validate()
	if err != nil {
		return err
	}

	err = p.buildManifest(ctx)
	if err != nil {
		return err
	}

	if o.GraphWidth != 1 {
		log.Info("🔨 Building graphs...")
		p.graphMD = buildGraphMD(p.body.Releases)
		log.Infof("show graph:\n%s", p.BuildGraphASCII(o.GraphWidth))
	}

	return nil
}
