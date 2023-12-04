package hooks

import (
	"context"
	"fmt"
	"os"

	"github.com/helmwave/helmwave/pkg/helper"
)

const (
	EnvReleaseUniqname = "HELMWAVE_LIFECYCLE_RELEASE_UNIQNAME"
	EnvLifecycleType   = "HELMWAVE_LIFECYCLE_TYPE"
)

func (h *hook) getCommandEnviron(ctx context.Context) []string {
	env := os.Environ()

	if uniq, exists := helper.ContextGetReleaseUniq(ctx); exists {
		env = addToEnviron(env, EnvReleaseUniqname, uniq.String())
	}

	if typ, exists := helper.ContextGetLifecycleType(ctx); exists {
		env = addToEnviron(env, EnvLifecycleType, typ)
	}

	return env
}

func addToEnviron(env []string, key, value string) []string {
	return append(env, fmt.Sprintf("%s=%s", key, value))
}
