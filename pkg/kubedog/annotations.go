package kubedog

const (
	// RootAnnoName is prefix for all kubedog annotations.
	RootAnnoName = "helmwave.dev/"

	// TrackTerminationModeAnnoName annotation allows to specify how to track resource.
	TrackTerminationModeAnnoName = RootAnnoName + "track-termination-mode"

	// FailModeAnnoName annotation specifies what to do after resource fails.
	FailModeAnnoName = RootAnnoName + "fail-mode"

	// FailuresAllowedPerReplicaAnnoName specifies how many times resource is allowed to fail.
	FailuresAllowedPerReplicaAnnoName = RootAnnoName + "failures-allowed-per-replica"

	// LogRegexAnnoName allows to set regexp for log lines.
	LogRegexAnnoName = RootAnnoName + "log-regex"

	// LogRegexForAnnoPrefix allows to set regexp for individual containers.
	LogRegexForAnnoPrefix = RootAnnoName + "log-regex-for-"

	// SkipLogsAnnoName allows to skip log streaming.
	SkipLogsAnnoName = RootAnnoName + "skip-logs"

	// SkipLogsForContainersAnnoName allows to skip log streaming for individual containers.
	SkipLogsForContainersAnnoName = RootAnnoName + "skip-logs-for-containers"

	// ShowLogsOnlyForContainersAnnoName allows to show logs only for specified containers.
	ShowLogsOnlyForContainersAnnoName = RootAnnoName + "show-logs-only-for-containers"

	// ShowLogsUntilAnnoName is unused.
	ShowLogsUntilAnnoName = RootAnnoName + "show-logs-until"

	// ShowEventsAnnoName enables streaming resource events.
	ShowEventsAnnoName = RootAnnoName + "show-service-messages"

	// ReplicasOnCreationAnnoName is unused.
	ReplicasOnCreationAnnoName = RootAnnoName + "replicas-on-creation"
)
