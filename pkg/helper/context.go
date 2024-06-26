package helper

import (
	"context"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
)

type contextReleaseUniqKey struct{}

func ContextWithReleaseUniq(ctx context.Context, name uniqname.UniqName) context.Context {
	return context.WithValue(ctx, contextReleaseUniqKey{}, name)
}

func ContextGetReleaseUniq(ctx context.Context) (uniqname.UniqName, bool) {
	v := ctx.Value(contextReleaseUniqKey{})

	switch v := v.(type) {
	case uniqname.UniqName:
		return v, true
	default:
		return uniqname.UniqName{}, false
	}
}

type contextLifecycleTypeKey struct{}

func ContextWithLifecycleType(ctx context.Context, typ string) context.Context {
	return context.WithValue(ctx, contextLifecycleTypeKey{}, typ)
}

func ContextGetLifecycleType(ctx context.Context) (string, bool) {
	v := ctx.Value(contextLifecycleTypeKey{})

	switch v := v.(type) {
	case string:
		return v, true
	default:
		return "", false
	}
}
