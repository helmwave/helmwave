package kubedog

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack"
)

// MakeSpecs creates *multitrack.MultitrackSpecs for Resource slice in provided namespace.
func MakeSpecs(m []Resource, ns string) (*multitrack.MultitrackSpecs, error) {
	specs := &multitrack.MultitrackSpecs{}

	for i := 0; i < len(m); i++ {
		r := &m[i]

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

// MakeMultiTrackSpec creates *multitrack.MultitrackSpec for current resource.
func (r *Resource) MakeMultiTrackSpec(ns string) (*multitrack.MultitrackSpec, error) {
	// Default spec
	spec := &multitrack.MultitrackSpec{
		ResourceName: r.Name,
		// Namespace:               r.Namespace,
		Namespace:               ns,
		LogRegexByContainerName: map[string]*regexp.Regexp{},
		TrackTerminationMode:    multitrack.WaitUntilResourceReady,
		FailMode:                multitrack.FailWholeDeployProcessImmediately,
		AllowFailuresCount:      new(int),
		FailureThresholdSeconds: new(int),
	}
	*spec.AllowFailuresCount = 1
	*spec.FailureThresholdSeconds = 0

	// Override by annotations
	for name, value := range r.Annotations {
		// invalid := fmt.Errorf("%s/%s annotation %s with invalid value %s", r.Name, r.Kind, name, value)
		var err error

		switch name {
		// Parse Value
		case SkipLogsAnnoName:
			err = r.handleAnnotationSkipLogs(value, spec)
		case ShowEventsAnnoName:
			err = r.handleAnnotationShowEvents(value, spec)
		case LogRegexAnnoName:
			err = r.handleAnnotationLogRegex(value, spec)
		case FailuresAllowedPerReplicaAnnoName:
			err = r.handleAnnotationFailuresAllowedPerReplica(value, spec)

		// Choose value
		case TrackTerminationModeAnnoName:
			err = r.handleAnnotationTrackTerminationMode(value, spec)
		case FailModeAnnoName:
			err = r.handleAnnotationFailMode(value, spec)

		// Parse array
		case SkipLogsForContainersAnnoName:
			err = r.handleAnnotationSkipLogsForContainers(value, spec)
		case ShowLogsOnlyForContainersAnnoName:
			err = r.handleAnnotationShowLogsOnlyForContainers(value, spec)

		default:
			//nolint:gocritic // keep switch in case of more prefix-based annotations in future
			switch {
			case strings.HasPrefix(name, LogRegexForAnnoPrefix):
				err = r.handleAnnotationLogRegexFor(name, value, spec)
			}
		}

		if err != nil {
			return nil, err
		}
	}

	return spec, nil
}

func (*Resource) handleAnnotationSkipLogs(value string, spec *multitrack.MultitrackSpec) error {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("failed to parse %s as boolean: %w", value, err)
	}
	spec.SkipLogs = v

	return nil
}

func (*Resource) handleAnnotationShowEvents(value string, spec *multitrack.MultitrackSpec) error {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("failed to parse %s as boolean: %w", value, err)
	}
	spec.ShowServiceMessages = v

	return nil
}

func (*Resource) handleAnnotationLogRegex(value string, spec *multitrack.MultitrackSpec) error {
	v, err := regexp.Compile(value)
	if err != nil {
		return fmt.Errorf("failed to compile %s as regexp: %w", value, err)
	}
	spec.LogRegex = v

	return nil
}

func (r *Resource) handleAnnotationFailuresAllowedPerReplica(value string, spec *multitrack.MultitrackSpec) error {
	if r.Kind == "Job" {
		return fmt.Errorf("%s is not supported for jobs", FailuresAllowedPerReplicaAnnoName)
	}

	v, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse %s as uint: %w", value, err)
	}

	replicas := 1
	if r.Spec.Replicas != nil {
		replicas = int(*r.Spec.Replicas)
	}

	*spec.AllowFailuresCount = int(v) * replicas

	return nil
}

func (*Resource) handleAnnotationTrackTerminationMode(value string, spec *multitrack.MultitrackSpec) error {
	v := multitrack.TrackTerminationMode(value)
	values := []multitrack.TrackTerminationMode{
		multitrack.WaitUntilResourceReady,
		multitrack.NonBlocking,
	}
	for _, mode := range values {
		if mode == v {
			spec.TrackTerminationMode = v

			return nil
		}
	}

	return fmt.Errorf("%s not found", v)
}

func (*Resource) handleAnnotationFailMode(value string, spec *multitrack.MultitrackSpec) error {
	v := multitrack.FailMode(value)
	values := []multitrack.FailMode{
		multitrack.IgnoreAndContinueDeployProcess,
		multitrack.FailWholeDeployProcessImmediately,
		multitrack.HopeUntilEndOfDeployProcess,
	}
	for _, mode := range values {
		if mode == v {
			spec.FailMode = v

			return nil
		}
	}

	return fmt.Errorf("%s not found", v)
}

func (*Resource) handleAnnotationSkipLogsForContainers(value string, spec *multitrack.MultitrackSpec) error {
	containers, err := splitContainers(value)
	if err != nil {
		return err
	}
	spec.SkipLogsForContainers = containers

	return nil
}

func (*Resource) handleAnnotationShowLogsOnlyForContainers(value string, spec *multitrack.MultitrackSpec) error {
	containers, err := splitContainers(value)
	if err != nil {
		return err
	}
	spec.ShowLogsOnlyForContainers = containers

	return nil
}

func (*Resource) handleAnnotationLogRegexFor(name, value string, spec *multitrack.MultitrackSpec) error {
	containerName := strings.TrimPrefix(name, LogRegexForAnnoPrefix)
	if containerName == "" {
		log.WithField("annotation", name).Error("annotation is invalid: cannot get container name")

		return nil
	}

	regexpValue, err := regexp.Compile(value)
	if err != nil {
		return fmt.Errorf("failed to parse %s as uint: %w", value, err)
	}

	spec.LogRegexByContainerName[containerName] = regexpValue

	return nil
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
