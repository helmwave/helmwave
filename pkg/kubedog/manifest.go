package kubedog

import (
	"bytes"
	"errors"
	"io"

	"github.com/goccy/go-yaml"
	log "github.com/sirupsen/logrus"
	meta1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Todo:  optimize?

// Resource is base structure for all k8s resources that have replicas.
// Used to parse out replicas count.
type Resource struct {
	Spec             `yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	meta1.TypeMeta   `yaml:",inline"`
	meta1.ObjectMeta `yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

// Spec is spec structure with replicas. Only replicas count is used.
type Spec struct {
	Replicas *uint32 `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`
}

// Parse parses YAML manifests of kubernetes resources and returns Resource slice.
//
//nolint:contextcheck,nolintlint
func Parse(yamlFile []byte) []Resource {
	var a []Resource

	r := bytes.NewReader(yamlFile)
	dec := yaml.NewDecoder(r)

	var t Resource
	var err error
	for ; !errors.Is(err, io.EOF); err = dec.Decode(&t) {
		if err != nil {
			log.WithError(err).Info("failed to parse resource manifest for kubedog")

			continue
		}
		a = append(a, t)
	}

	return a
}
