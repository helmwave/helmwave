package kubedog

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	log "github.com/sirupsen/logrus"
	meta1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

// Todo:  optimize?

// Resource is base structure for all k8s resources that have replicas.
// Used to parse out replicas count.
type Resource struct {
	Spec             `yaml:"spec,omitempty" json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	meta1.TypeMeta   `yaml:",inline" json:",inline"`
	meta1.ObjectMeta `yaml:"metadata,omitempty" json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

// DeepCopyObject is required to implement runtime.Object interface. It doesn't actually do anything, don't use.
func (r *Resource) DeepCopyObject() runtime.Object {
	return r
}

// Spec is spec structure with replicas. Only replicas count is used.
type Spec struct {
	Replicas *uint32 `yaml:"replicas,omitempty" json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`
}

// Parse parses YAML manifests of kubernetes resources and returns Resource slice.
func Parse(yamlFile []byte) []Resource {
	var a []Resource

	r := bytes.NewReader(yamlFile)
	dec := yaml.NewYAMLReader(bufio.NewReader(r))
	d := scheme.Codecs.UniversalDeserializer()
	var doc []byte
	var err error
	for ; !errors.Is(err, io.EOF); doc, err = dec.Read() {
		if err != nil {
			log.WithError(err).Info("failed to parse resource manifest for kubedog")

			continue
		}

		var t Resource
		_, _, err := d.Decode(doc, nil, &t)
		if err != nil {
			log.WithError(err).Info("failed to parse resource manifest for kubedog")

			continue
		}

		a = append(a, t)
	}

	return a
}
