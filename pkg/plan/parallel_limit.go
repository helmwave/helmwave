package plan

import (
	"context"

	"github.com/helmwave/helmwave/pkg/clictx"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) ParallelLimiter(ctx context.Context) int {
	return ParallelLimit(ctx, p.body.Releases)
}

func ParallelLimit(ctx context.Context, releases release.Configs) int {
	limit, ok := clictx.GetFlagFromContext(ctx, "parallel-limit").(int)
	if !ok {
		limit = 0
	}
	if limit == 0 {
		limit = len(releases)
	}

	const msg = "Releases limited parallelization"
	if limit == len(releases) {
		log.WithField("limit", limit).Debug(msg)
	} else {
		log.WithField("limit", limit).Info(msg)
	}

	return limit
}

// func (p *Plan) parallelWorker(ctx context.Context, fn func(any)) {
//	ch := p.Graph().Run()
//	wg := parallel.NewWaitGroup()
//	limiter := p.ParallelLimiter(ctx)
//	wg.Add(limiter)
//
//	for range limiter {
//		go func() {
//			defer wg.Done()
//			for n := range ch {
//				fn()
//			}
//		}()
//	}
//}
