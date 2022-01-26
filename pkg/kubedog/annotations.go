package kubedog

import "github.com/helmwave/helmwave/pkg/helper"

const (
	// TrackTerminationModeAnnoName annotation allows to specify how to track resource.
	TrackTerminationModeAnnoName = helper.RootAnnoName + "track-termination-mode"

	// FailModeAnnoName annotation specifies what to do after resource fails.
	FailModeAnnoName = helper.RootAnnoName + "fail-mode"

	// FailuresAllowedPerReplicaAnnoName specifies how many times resource is allowed to fail.
	FailuresAllowedPerReplicaAnnoName = helper.RootAnnoName + "failures-allowed-per-replica"

	// LogRegexAnnoName allows to set regexp for log lines.
	LogRegexAnnoName = helper.RootAnnoName + "log-regex"

	// LogRegexForAnnoPrefix allows to set regexp for individual containers.
	LogRegexForAnnoPrefix = helper.RootAnnoName + "log-regex-for-"

	// SkipLogsAnnoName allows to skip log streaming.
	SkipLogsAnnoName = helper.RootAnnoName + "skip-logs"

	// SkipLogsForContainersAnnoName allows to skip log streaming for individual containers.
	SkipLogsForContainersAnnoName = helper.RootAnnoName + "skip-logs-for-containers"

	// ShowLogsOnlyForContainersAnnoName allows to show logs only for specified containers.
	ShowLogsOnlyForContainersAnnoName = helper.RootAnnoName + "show-logs-only-for-containers"

	// ShowLogsUntilAnnoName is unused.
	ShowLogsUntilAnnoName = helper.RootAnnoName + "show-logs-until"

	// ShowEventsAnnoName enables streaming resource events.
	ShowEventsAnnoName = helper.RootAnnoName + "show-service-messages"

	// ReplicasOnCreationAnnoName is unused.
	ReplicasOnCreationAnnoName = helper.RootAnnoName + "replicas-on-creation"
)
