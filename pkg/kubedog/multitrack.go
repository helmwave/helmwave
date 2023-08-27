package kubedog

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"github.com/werf/kubedog/pkg/tracker/resid"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack/generic"
	"golang.org/x/exp/slices"
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
		case SkipLogsAnnoName, OldSkipLogsAnnoName:
			err = r.handleAnnotationSkipLogs(value, spec)
		case ShowEventsAnnoName, OldShowEventsAnnoName:
			err = r.handleAnnotationShowEvents(value, spec)
		case LogRegexAnnoName, OldLogRegexAnnoName:
			err = r.handleAnnotationLogRegex(value, spec)
		case FailuresAllowedPerReplicaAnnoName, OldFailuresAllowedPerReplicaAnnoName:
			err = r.handleAnnotationFailuresAllowedPerReplica(value, spec)

		// Choose value
		case TrackTerminationModeAnnoName, OldTrackTerminationModeAnnoName:
			err = r.handleAnnotationTrackTerminationMode(name, value, spec)
		case FailModeAnnoName, OldFailModeAnnoName:
			err = r.handleAnnotationFailMode(name, value, spec)

		// Parse array
		case SkipLogsForContainersAnnoName, OldSkipLogsForContainersAnnoName:
			err = r.handleAnnotationSkipLogsForContainers(name, value, spec)
		case ShowLogsOnlyForContainersAnnoName, OldShowLogsOnlyForContainersAnnoName:
			err = r.handleAnnotationShowLogsOnlyForContainers(name, value, spec)

		default:
			switch {
			case strings.HasPrefix(name, LogRegexForAnnoPrefix), strings.HasPrefix(name, OldLogRegexForAnnoPrefix):
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
		return NewParseError("boolean", value, err)
	}
	spec.SkipLogs = v

	return nil
}

func (*Resource) handleAnnotationShowEvents(value string, spec *multitrack.MultitrackSpec) error {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return NewParseError("boolean", value, err)
	}
	spec.ShowServiceMessages = v

	return nil
}

func (*Resource) handleAnnotationLogRegex(value string, spec *multitrack.MultitrackSpec) error {
	v, err := regexp.Compile(value)
	if err != nil {
		return NewParseError("regexp", value, err)
	}
	spec.LogRegex = v

	return nil
}

func (r *Resource) handleAnnotationFailuresAllowedPerReplica(value string, spec *multitrack.MultitrackSpec) error {
	v, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return NewParseError("uint", value, err)
	}

	replicas := 1
	if r.Spec.Replicas != nil {
		replicas = int(*r.Spec.Replicas)
	}

	*spec.AllowFailuresCount = int(v) * replicas

	return nil
}

func (*Resource) handleAnnotationTrackTerminationMode(anno, value string, spec *multitrack.MultitrackSpec) error {
	v := multitrack.TrackTerminationMode(value)
	values := []multitrack.TrackTerminationMode{
		multitrack.WaitUntilResourceReady,
		multitrack.NonBlocking,
	}

	if slices.Contains(values, v) {
		spec.TrackTerminationMode = v

		return nil
	}

	return NewInvalidValueError(anno, v, values)
}

func (*Resource) handleAnnotationFailMode(anno, value string, spec *multitrack.MultitrackSpec) error {
	v := multitrack.FailMode(value)
	values := []multitrack.FailMode{
		multitrack.IgnoreAndContinueDeployProcess,
		multitrack.FailWholeDeployProcessImmediately,
		multitrack.HopeUntilEndOfDeployProcess,
	}

	if slices.Contains(values, v) {
		spec.FailMode = v

		return nil
	}

	return NewInvalidValueError(anno, v, values)
}

func (*Resource) handleAnnotationSkipLogsForContainers(name, value string, spec *multitrack.MultitrackSpec) error {
	containers, err := splitContainers(name, value)
	if err != nil {
		return err
	}
	spec.SkipLogsForContainers = containers

	return nil
}

func (*Resource) handleAnnotationShowLogsOnlyForContainers(name, value string, spec *multitrack.MultitrackSpec) error {
	containers, err := splitContainers(name, value)
	if err != nil {
		return err
	}
	spec.ShowLogsOnlyForContainers = containers

	return nil
}

func (*Resource) handleAnnotationLogRegexFor(name, value string, spec *multitrack.MultitrackSpec) error {
	var containerName string
	switch {
	case strings.HasPrefix(name, LogRegexForAnnoPrefix):
		containerName = strings.TrimPrefix(name, LogRegexForAnnoPrefix)
	case strings.HasPrefix(name, OldLogRegexForAnnoPrefix):
		containerName = strings.TrimPrefix(name, OldLogRegexForAnnoPrefix)
	}

	if containerName == "" {
		log.WithField("annotation", name).Error("annotation is invalid: can't get container name")

		return nil
	}

	regexpValue, err := regexp.Compile(value)
	if err != nil {
		return NewParseError("uint", value, err)
	}

	spec.LogRegexByContainerName[containerName] = regexpValue

	return nil
}

func splitContainers(name, value string) (containers []string, err error) {
	for _, v := range strings.Split(value, ",") {
		container := strings.TrimSpace(v)
		if container == "" {
			return nil, NewEmptyContainerNameError(name, value)
		}

		containers = append(containers, container)
	}

	return containers, err
}

func isIgnoredGenericGK(gk schema.GroupKind) bool {
	l := strings.ToLower(gk.String())

	return helper.Contains(l, ignoredGenericGK)
}
