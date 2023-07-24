package kubedog

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"github.com/werf/kubedog/pkg/tracker/resid"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack/generic"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var ignoredGenericGK = []string{
	"componentstatus",
	"namespace",
	"node",
	"persistentvolume",
	"mutatingwebhookconfiguration.admissionregistration.k8s.io",
	"validatingwebhookconfiguration.admissionregistration.k8s.io",
	"customresourcedefinition.apiextensions.k8s.io",
	"apiservice.apiregistration.k8s.io",
	"tokenreview.authentication.k8s.io",
	"selfsubjectaccessreview.authorization.k8s.io",
	"selfsubjectrulesreview.authorization.k8s.io",
	"subjectaccessreview.authorization.k8s.io",
	"certificatesigningrequest.certificates.k8s.io",
	"flowschema.flowcontrol.apiserver.k8s.io",
	"prioritylevelconfiguration.flowcontrol.apiserver.k8s.io",
	"ingressclass.networking.k8s.io",
	"runtimeclass.node.k8s.io",
	"clusterrolebinding.rbac.authorization.k8s.io",
	"clusterrole.rbac.authorization.k8s.io",
	"priorityclass.scheduling.k8s.io",
	"csidriver.storage.k8s.io",
	"csinode.storage.k8s.io",
	"storageclass.storage.k8s.io",
	"volumeattachment.storage.k8s.io",
}

// MakeSpecs creates *multitrack.MultitrackSpecs for Resource slice in provided namespace.
func MakeSpecs(m []Resource, ns string, trackGeneric bool) (*multitrack.MultitrackSpecs, error) {
	specs := &multitrack.MultitrackSpecs{}

	for i := 0; i < len(m); i++ {
		r := &m[i]

		spec, err := r.MakeMultiTrackSpec(ns)
		if err != nil {
			return nil, err
		}

		switch r.Kind {
		case "Deployment":
			specs.Deployments = append(specs.Deployments, *spec)
		case "StatefulSet":
			specs.StatefulSets = append(specs.StatefulSets, *spec)
		case "DaemonSet":
			specs.DaemonSets = append(specs.DaemonSets, *spec)
		case "Job":
			specs.Jobs = append(specs.Jobs, *spec)
		case "Canary":
			specs.Canaries = append(specs.Canaries, *spec)
		case "": // probably some empty manifest due to templating, just skipping it
		default:
			if !trackGeneric {
				continue
			}

			// skipping some common cluster-wide resources because they are not supported by kubedog
			if isIgnoredGenericGK(r.GroupVersionKind().GroupKind()) {
				continue
			}
			s := &generic.Spec{
				ResourceID: &resid.ResourceID{
					Name:             spec.ResourceName,
					Namespace:        spec.Namespace,
					GroupVersionKind: r.GroupVersionKind(),
				},
				Timeout:              0,
				NoActivityTimeout:    nil,
				TrackTerminationMode: generic.TrackTerminationMode(spec.TrackTerminationMode),
				FailMode:             generic.FailMode(spec.FailMode),
				AllowFailuresCount:   spec.AllowFailuresCount,
				ShowServiceMessages:  spec.ShowServiceMessages,
				HideEvents:           false,
				StatusProgressPeriod: 0,
			}
			err := s.Init()
			if err != nil {
				log.WithError(err).
					WithField("resource name", spec.ResourceName).
					WithField("resource type", r.GroupVersionKind().String()).
					WithField("resource manifest", r.Spec).
					Warn("failed to create watcher for resource, skipping the resource")

				continue
			}
			specs.Generics = append(specs.Generics, s)
		}
	}

	return specs, nil
}

// MakeMultiTrackSpec creates *multitrack.MultitrackSpec for current resource.
func (r *Resource) MakeMultiTrackSpec(ns string) (*multitrack.MultitrackSpec, error) {
	// Default spec
	spec := &multitrack.MultitrackSpec{
		ResourceName:            r.Name,
		Namespace:               ns,
		LogRegexByContainerName: map[string]*regexp.Regexp{},
		TrackTerminationMode:    multitrack.WaitUntilResourceReady,
		FailMode:                multitrack.FailWholeDeployProcessImmediately,
		AllowFailuresCount:      new(int),
		FailureThresholdSeconds: new(int),
	}
	*spec.AllowFailuresCount = 0
	*spec.FailureThresholdSeconds = 0

	// Override by annotations
	for name, value := range r.Annotations {
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
			switch { //nolint:gocritic // keep switch in case of more prefix-based annotations in future
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
		log.WithField("annotation", name).Error("annotation is invalid: can't get container name")

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

func isIgnoredGenericGK(gk schema.GroupKind) bool {
	l := strings.ToLower(gk.String())

	return helper.Contains(l, ignoredGenericGK)
}
