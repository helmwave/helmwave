package plan

import (
	"context"
	"github.com/helmwave/helmwave/pkg/clictx"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) getParallelLimit(ctx context.Context) int {
	return getParallelLimit(ctx, p.body.Releases)
}

func getParallelLimit(ctx context.Context, releases release.Configs) int {
	parallelLimit, ok := clictx.GetFlagFromContext(ctx, "parallel-limit").(int)
	if !ok {
		parallelLimit = 0
	}
	if parallelLimit == 0 {
		parallelLimit = len(releases)
	}

	const msg = "Releases limited parallelization"
	if parallelLimit == len(releases) {
		log.WithField("limit", parallelLimit).Debug(msg)
	} else {
		log.WithField("limit", parallelLimit).Info(msg)
	}

	return parallelLimit
}
