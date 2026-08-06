package release

import (
	"github.com/invopop/jsonschema"
	"helm.sh/helm/v4/pkg/kube"
)

// WaitStrategy tells helm how to wait for the resources of a release:
// watcher, legacy or hookOnly. Anything else is rejected at validation.
type WaitStrategy string

const (
	// WaitStrategyWatcher waits for every resource using kubernetes watches.
	WaitStrategyWatcher = WaitStrategy(kube.StatusWatcherStrategy)

	// WaitStrategyLegacy waits by polling.
	WaitStrategyLegacy = WaitStrategy(kube.LegacyStrategy)

	// WaitStrategyHookOnly waits only for hooks, not for the chart's own resources.
	// The default when `wait` is not set at all.
	WaitStrategyHookOnly = WaitStrategy(kube.HookOnlyStrategy)
)

// Helm returns the helm wait strategy. helm has no usable zero value here — it refuses to build a
// waiter for an empty strategy — so an unset strategy becomes the same default helm's own CLI uses.
func (s WaitStrategy) Helm() kube.WaitStrategy {
	if s == "" {
		return kube.HookOnlyStrategy
	}

	return kube.WaitStrategy(s)
}

// Enabled reports whether the release waits for its own resources, and not just for its hooks.
func (s WaitStrategy) Enabled() bool {
	return s.Helm() != kube.HookOnlyStrategy
}

// Validate checks that the strategy is one helm knows.
func (s WaitStrategy) Validate() error {
	switch s {
	case "", WaitStrategyWatcher, WaitStrategyLegacy, WaitStrategyHookOnly:
		return nil
	default:
		return NewInvalidWaitStrategyError(string(s))
	}
}

func (WaitStrategy) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: jsonSchemaString,
		Enum: []any{
			WaitStrategyWatcher,
			WaitStrategyLegacy,
			WaitStrategyHookOnly,
			"",
		},
	}
}
