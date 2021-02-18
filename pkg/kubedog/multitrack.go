package kubedog

import (
	"fmt"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack"
	"regexp"
	"strconv"
	"strings"
)

func MakeSpecs(m []Resource, ns string) (*multitrack.MultitrackSpecs, error) {
	specs := &multitrack.MultitrackSpecs{}

	for _, r := range m {
		switch r.Kind {
		case "Deployment":
			s, err := r.MakeMultiTrackSpec(ns)
			if err != nil {
				return nil, err
			}
			specs.Deployments = append(specs.Deployments, *s)
		case "StatefulSet":
			s, err := r.MakeMultiTrackSpec(ns)
			if err != nil {
				return nil, err
			}
			specs.StatefulSets = append(specs.StatefulSets, *s)
		case "Job":
			s, err := r.MakeMultiTrackSpec(ns)
			if err != nil {
				return nil, err
			}
			specs.Jobs = append(specs.Jobs, *s)
		case "DaemonSet":
			s, err := r.MakeMultiTrackSpec(ns)
			if err != nil {
				return nil, err
			}
			specs.DaemonSets = append(specs.DaemonSets, *s)
		}

	}

	return specs, nil
}

// BolgenOS on max
func (r *Resource) MakeMultiTrackSpec(ns string) (*multitrack.MultitrackSpec, error) {
	// Default spec
	spec := &multitrack.MultitrackSpec{
		ResourceName: r.Name,
		//Namespace:               r.Namespace,
		Namespace:               ns,
		LogRegexByContainerName: map[string]*regexp.Regexp{},
		TrackTerminationMode:    multitrack.WaitUntilResourceReady,
		FailMode:                multitrack.FailWholeDeployProcessImmediately,
	}
	spec.AllowFailuresCount = new(int)
	*spec.AllowFailuresCount = 1
	spec.FailureThresholdSeconds = new(int)
	*spec.FailureThresholdSeconds = 0

	// Override by annotations
loop:
	for name, value := range r.Annotations {
		// invalid := fmt.Errorf("%s/%s annotation %s with invalid value %s", r.Name, r.Kind, name, value)

		switch name {

		// Parse Value
		case SkipLogsAnnoName:
			v, err := strconv.ParseBool(value)
			if err != nil {
				return nil, err
			}
			spec.SkipLogs = v
		case ShowEventsAnnoName:
			v, err := strconv.ParseBool(value)
			if err != nil {
				return nil, err
			}
			spec.ShowServiceMessages = v
		case LogRegexAnnoName:
			v, err := regexp.Compile(value)
			if err != nil {
				return nil, err
			}
			spec.LogRegex = v
		case FailuresAllowedPerReplicaAnnoName:
			if r.Kind == "Job" {
				return nil, fmt.Errorf("%s does not support for job", name)
			}

			v, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, err
			}

			if r.Spec.Replicas == nil {
				*r.Spec.Replicas = 1
			}

			*spec.AllowFailuresCount = int(v) * int(*r.Spec.Replicas)

		// Chose value
		case TrackTerminationModeAnnoName:
			v := multitrack.TrackTerminationMode(value)
			values := []multitrack.TrackTerminationMode{
				multitrack.WaitUntilResourceReady,
				multitrack.NonBlocking,
			}
			for _, mode := range values {
				if mode == v {
					spec.TrackTerminationMode = v
					continue loop
				}
			}
		case FailModeAnnoName:
			v := multitrack.FailMode(value)
			values := []multitrack.FailMode{
				multitrack.IgnoreAndContinueDeployProcess,
				multitrack.FailWholeDeployProcessImmediately,
				multitrack.HopeUntilEndOfDeployProcess,
			}
			for _, mode := range values {
				if mode == v {
					spec.FailMode = v
					continue loop
				}
			}

			return nil, fmt.Errorf("%s not found", v)

		// Parse array
		case SkipLogsForContainersAnnoName:
			containers, err := splitContainers(value)
			if err != nil {
				return nil, err
			}
			spec.SkipLogsForContainers = containers
		case ShowLogsOnlyForContainers:
			containers, err := splitContainers(value)
			if err != nil {
				return nil, err
			}
			spec.ShowLogsOnlyForContainers = containers

		default:
			if strings.HasPrefix(name, LogRegexForAnnoPrefix) {
				if containerName := strings.TrimPrefix(name, LogRegexForAnnoPrefix); containerName != "" {
					regexpValue, err := regexp.Compile(value)
					if err != nil {
						return nil, err
					}

					spec.LogRegexByContainerName[containerName] = regexpValue
				}
			}
		}
	}

	return spec, nil
}

func splitContainers(annoValue string) (containers []string, err error) {
	for _, v := range strings.Split(annoValue, ",") {
		container := strings.TrimSpace(v)
		if container == "" {
			return nil, fmt.Errorf("%s: containers names separated by comma expected", annoValue)
		}

		containers = append(containers, container)
	}

	return containers, err
}
