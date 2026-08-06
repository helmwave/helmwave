package tracker_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/tracker"
	"github.com/stretchr/testify/suite"
)

const testManifest = `
apiVersion: flagger.app/v1beta1
kind: Canary
metadata:
  name: blabla
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: blabla
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: blabla
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: bla
---
apiVersion: batch/v1
kind: Job
metadata:
  name: blabla
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: blabla
  namespace: custom
---
`

type ManifestTestSuite struct {
	suite.Suite
}

func TestManifestTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ManifestTestSuite))
}

func (s *ManifestTestSuite) TestWorkloadsOnly() {
	ids := tracker.ObjectsFromManifest(testManifest, "ns", false)

	s.Require().Len(ids, 4)
	for i := range ids {
		id := &ids[i]
		s.Contains([]string{"Deployment", "StatefulSet", "DaemonSet", "Job"}, id.GroupKind.Kind)
	}
}

func (s *ManifestTestSuite) TestTrackAll() {
	ids := tracker.ObjectsFromManifest(testManifest, "ns", true)

	// ServiceAccount is ignored, the Canary CR is tracked generically.
	s.Require().Len(ids, 5)

	kinds := make([]string, 0, len(ids))
	for i := range ids {
		id := &ids[i]
		kinds = append(kinds, id.GroupKind.Kind)
	}
	s.Contains(kinds, "Canary")
	s.NotContains(kinds, "ServiceAccount")
}

func (s *ManifestTestSuite) TestNamespaceDefaulting() {
	ids := tracker.ObjectsFromManifest(testManifest, "default-ns", false)

	for i := range ids {
		id := &ids[i]
		if id.GroupKind.Kind == "StatefulSet" {
			s.Equal("custom", id.Namespace)
		} else {
			s.Equal("default-ns", id.Namespace)
		}
	}
}

func (s *ManifestTestSuite) TestDeduplicate() {
	manifest := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: same
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: same
`
	ids := tracker.ObjectsFromManifest(manifest, "ns", false)

	s.Require().Len(ids, 1)
}

func (s *ManifestTestSuite) TestSkipsGarbage() {
	manifest := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ok
---
kind: NoName
metadata: {}
---
{{ not yaml at all
`
	ids := tracker.ObjectsFromManifest(manifest, "ns", false)

	s.Require().Len(ids, 1)
	s.Equal("ok", ids[0].Name)
}

func (s *ManifestTestSuite) TestEmpty() {
	s.Empty(tracker.ObjectsFromManifest("", "ns", true))
}
