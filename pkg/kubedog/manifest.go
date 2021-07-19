package kubedog

import (
	"bytes"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Todo:  optimize?

type Resource struct {
	metav1.TypeMeta   `yaml:",inline"`
	metav1.ObjectMeta `yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              `yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type Spec struct {
	Replicas *uint32 `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`
}

func Parse(yamlFile []byte) []Resource {
	var a []Resource

	r := bytes.NewReader(yamlFile)
	dec := yaml.NewDecoder(r)

	var t Resource
	for dec.Decode(&t) == nil {
		a = append(a, t)
	}

	return a
}
