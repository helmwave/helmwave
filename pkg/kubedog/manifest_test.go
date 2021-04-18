package kubedog

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ManifestTestSuite struct {
	suite.Suite
}

func (s *ManifestTestSuite) TestMakeManifest() {
	data := `
apiVersion: blalba/v1
kind: Blablaployment
metadata:
  name: "123"
  annotations:
    bla: "blabla"
spec:
  replicas: 123
---
apiVersion: blabla/v1
kind: BlablaSet
metadata:
  name: "123"
  annotations:
    bla: "blabla"
spec:
  replicas: 123
`
	resources := MakeManifest([]byte(data))
	s.Require().Len(resources, 2)

	s.Equal("Blablaployment", resources[0].Kind)
	s.Equal("BlablaSet", resources[1].Kind)

	s.Equal("123", resources[0].ObjectMeta.Name)
	s.Equal("123", resources[1].ObjectMeta.Name)

	s.Require().NotNil(resources[0].ObjectMeta.Annotations)
	s.Require().NotNil(resources[0].ObjectMeta.Annotations)

	s.Equal("blabla", resources[0].ObjectMeta.Annotations["bla"])
	s.Equal("blabla", resources[1].ObjectMeta.Annotations["bla"])

	s.Require().NotNil(resources[0].Spec.Replicas)
	s.Require().NotNil(resources[1].Spec.Replicas)

	s.Equal(uint32(123), *resources[0].Spec.Replicas)
	s.Equal(uint32(123), *resources[1].Spec.Replicas)
}

func TestManifestTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ManifestTestSuite))
}
