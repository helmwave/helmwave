package tracker

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/fluxcd/cli-utils/pkg/object"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
)

const (
	groupApps      = "apps"
	groupBatch     = "batch"
	kindDeployment = "Deployment"
	kindJob        = "Job"
	groupAuthz     = "authorization.k8s.io"
	groupRBAC      = "rbac.authorization.k8s.io"
	groupStorage   = "storage.k8s.io"
)

// workloadGK is the set of kinds tracked by default.
//
//nolint:gochecknoglobals // static lookup table
var workloadGK = map[schema.GroupKind]bool{
	{Group: groupApps, Kind: kindDeployment}: true,
	{Group: groupApps, Kind: "StatefulSet"}:  true,
	{Group: groupApps, Kind: "DaemonSet"}:    true,
	{Group: groupBatch, Kind: kindJob}:       true,
}

// ignoredGK are cluster-wide or status-less kinds that only produce noise when tracked.
//
//nolint:gochecknoglobals // static lookup table
var ignoredGK = map[schema.GroupKind]bool{
	{Group: "", Kind: "ComponentStatus"}:                                            true,
	{Group: "", Kind: "Namespace"}:                                                  true,
	{Group: "", Kind: "Node"}:                                                       true,
	{Group: "", Kind: "PersistentVolume"}:                                           true,
	{Group: "", Kind: "Secret"}:                                                     true,
	{Group: "", Kind: "ConfigMap"}:                                                  true,
	{Group: "", Kind: "ServiceAccount"}:                                             true,
	{Group: "admissionregistration.k8s.io", Kind: "MutatingWebhookConfiguration"}:   true,
	{Group: "admissionregistration.k8s.io", Kind: "ValidatingWebhookConfiguration"}: true,
	{Group: "apiextensions.k8s.io", Kind: "CustomResourceDefinition"}:               true,
	{Group: "apiregistration.k8s.io", Kind: "APIService"}:                           true,
	{Group: "authentication.k8s.io", Kind: "TokenReview"}:                           true,
	{Group: groupAuthz, Kind: "SelfSubjectAccessReview"}:                            true,
	{Group: groupAuthz, Kind: "SelfSubjectRulesReview"}:                             true,
	{Group: groupAuthz, Kind: "SubjectAccessReview"}:                                true,
	{Group: "certificates.k8s.io", Kind: "CertificateSigningRequest"}:               true,
	{Group: "flowcontrol.apiserver.k8s.io", Kind: "FlowSchema"}:                     true,
	{Group: "flowcontrol.apiserver.k8s.io", Kind: "PriorityLevelConfiguration"}:     true,
	{Group: "networking.k8s.io", Kind: "IngressClass"}:                              true,
	{Group: "node.k8s.io", Kind: "RuntimeClass"}:                                    true,
	{Group: groupRBAC, Kind: "ClusterRoleBinding"}:                                  true,
	{Group: groupRBAC, Kind: "ClusterRole"}:                                         true,
	{Group: groupRBAC, Kind: "RoleBinding"}:                                         true,
	{Group: groupRBAC, Kind: "Role"}:                                                true,
	{Group: "scheduling.k8s.io", Kind: "PriorityClass"}:                             true,
	{Group: groupStorage, Kind: "CSIDriver"}:                                        true,
	{Group: groupStorage, Kind: "CSINode"}:                                          true,
	{Group: groupStorage, Kind: "StorageClass"}:                                     true,
	{Group: groupStorage, Kind: "VolumeAttachment"}:                                 true,
}

type manifestDoc struct {
	APIVersion string `yaml:"apiVersion"` //nolint:tagliatelle // the kubernetes field is camelCase
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
}

// ObjectsFromManifest parses multi-document YAML manifests and returns the objects to track.
// Namespace-less objects get defaultNamespace. Without trackAll only workload kinds are
// returned; with it, everything except ignoredGK.
//
// Documents are split before parsing so a malformed one only loses itself: yaml.v3's stream
// decoder cannot recover after a syntax error and would return it forever.
func ObjectsFromManifest(manifest, defaultNamespace string, trackAll bool) object.ObjMetadataSet {
	var ids object.ObjMetadataSet
	seen := map[object.ObjMetadata]bool{}

	r := utilyaml.NewYAMLReader(bufio.NewReader(strings.NewReader(manifest)))
	for {
		raw, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.WithError(err).Debug("failed to split manifests for tracking")

			break
		}

		var doc manifestDoc
		if err := yaml.Unmarshal(raw, &doc); err != nil {
			log.WithError(err).Debug("failed to parse a manifest document for tracking, skipping it")

			continue
		}

		id, ok := objectFromDoc(&doc, defaultNamespace, trackAll)
		if !ok || seen[id] {
			continue
		}
		seen[id] = true
		ids = append(ids, id)
	}

	return ids
}

func objectFromDoc(doc *manifestDoc, defaultNamespace string, trackAll bool) (object.ObjMetadata, bool) {
	if doc.Kind == "" || doc.Metadata.Name == "" {
		return object.ObjMetadata{}, false
	}

	gv, err := schema.ParseGroupVersion(doc.APIVersion)
	if err != nil {
		log.WithError(err).Debug("failed to parse apiVersion for tracking, skipping the document")

		return object.ObjMetadata{}, false
	}
	gk := gv.WithKind(doc.Kind).GroupKind()

	if trackAll {
		if ignoredGK[gk] {
			return object.ObjMetadata{}, false
		}
	} else if !workloadGK[gk] {
		return object.ObjMetadata{}, false
	}

	ns := doc.Metadata.Namespace
	if ns == "" {
		ns = defaultNamespace
	}

	return object.ObjMetadata{
		Namespace: ns,
		Name:      doc.Metadata.Name,
		GroupKind: gk,
	}, true
}
