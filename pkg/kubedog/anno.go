package kubedog

const (
	// Todo buy the domain
	RootAnnoName = "helmwave.dev/"

	TrackTerminationModeAnnoName      = RootAnnoName + "track-termination-mode"
	FailModeAnnoName                  = RootAnnoName + "fail-mode"
	FailuresAllowedPerReplicaAnnoName = RootAnnoName + "failures-allowed-per-replica"
	LogRegexAnnoName                  = RootAnnoName + "log-regex"
	LogRegexForAnnoPrefix             = RootAnnoName + "log-regex-for-"
	SkipLogsAnnoName                  = RootAnnoName + "skip-logs"
	SkipLogsForContainersAnnoName     = RootAnnoName + "skip-logs-for-containers"
	ShowLogsOnlyForContainers         = RootAnnoName + "show-logs-only-for-containers"
	// ShowLogsUntilAnnoName         	= RootAnnoName+"show-logs-until"
	ShowEventsAnnoName = RootAnnoName + "show-service-messages"
	// ReplicasOnCreationAnnoName 		= RootAnnoName+"replicas-on-creation"
)
