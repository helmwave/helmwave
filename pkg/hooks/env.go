package hooks

import (
	"context"
	"fmt"
	"os"

	"github.com/helmwave/helmwave/pkg/helper"
)

const (
	ENV_RELEASE_UNIQNAME = "HELMWAVE_LIFECYCLE_RELEASE_UNIQNAME"
	ENV_LIFECYCLE_TYPE   = "HELMWAVE_LIFECYCLE_TYPE"
)

func (h *hook) getCommandEnviron(ctx context.Context) []string {
	env := os.Environ()

	if uniq, exists := helper.ContextGetReleaseUniq(ctx); exists {
		env = addToEnviron(env, ENV_RELEASE_UNIQNAME, uniq.String())
	}

	if typ, exists := helper.ContextGetLifecycleType(ctx); exists {
		env = addToEnviron(env, ENV_LIFECYCLE_TYPE, typ)
	}

	return env
}

func addToEnviron(env []string, key, value string) []string {
	return append(env, fmt.Sprintf("%s=%s", key, value))
}
