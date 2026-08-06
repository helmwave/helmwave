package plan

import (
	"errors"
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/tracker"
	"github.com/stretchr/testify/suite"
	helmRelease "helm.sh/helm/v4/pkg/release/v1"
)

type TrackerTestSuite struct {
	suite.Suite
}

func TestTrackerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TrackerTestSuite))
}

const trackerTestManifest = `
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
---
`

func (s *TrackerTestSuite) TestNoReleases() {
	p := New("")
	p.NewBody()

	ids, _, err := p.trackerObjects(&tracker.Config{}, nil)

	s.Require().NoError(err)
	s.Require().Empty(ids)
}

func (s *TrackerTestSuite) TestCallsManifestFunction() {
	p := New("")
	p.NewBody()

	rel := NewMockReleaseConfig(s.T())
	rel.On("KubeContext").Return("")
	p.SetReleases(rel)

	s.Require().Panics(func() {
		_, _, _ = p.trackerObjects(&tracker.Config{}, nil)
	})
}

func (s *TrackerTestSuite) TestSyncObjects() {
	p := New("")
	p.NewBody()

	relName := "bla"
	relNS := "blabla"
	kubecontext := "blacontext"
	u, _ := uniqname.New(relName, relNS, "")

	p.manifests[u] = trackerTestManifest

	mockedRelease := NewMockReleaseConfig(s.T())
	mockedRelease.On("KubeContext").Return(kubecontext)
	mockedRelease.On("Uniq").Return(u)
	mockedRelease.On("Namespace").Return(relNS)
	p.SetReleases(mockedRelease)

	ids, context, err := p.trackerSyncObjects(&tracker.Config{TrackGeneric: true})

	s.Require().NoError(err)

	// ServiceAccount is ignored, everything else is tracked.
	s.Require().Len(ids, 5)
	for i := range ids {
		id := &ids[i]
		s.Equal(relNS, id.Namespace)
	}

	s.Require().Equal(kubecontext, context)

	mockedRelease.AssertExpectations(s.T())
}

func (s *TrackerTestSuite) TestRollbackObjects() {
	p := New("")
	p.NewBody()

	relName := "bla"
	relNS := "blabla"
	kubecontext := "blacontext"
	version := 666
	u, _ := uniqname.New(relName, relNS, "")

	mockedRelease := NewMockReleaseConfig(s.T())
	mockedRelease.On("KubeContext").Return(kubecontext)
	mockedRelease.On("Uniq").Return(u)
	mockedRelease.On("Namespace").Return(relNS)
	mockedRelease.On("Get", version).Return(&helmRelease.Release{Manifest: trackerTestManifest}, nil)
	p.SetReleases(mockedRelease)

	ids, context, err := p.trackerRollbackObjects(version, &tracker.Config{TrackGeneric: true})

	s.Require().NoError(err)
	s.Require().Len(ids, 5)
	s.Require().Equal(kubecontext, context)

	mockedRelease.AssertExpectations(s.T())
}

func (s *TrackerTestSuite) TestRollbackObjectsGetError() {
	p := New("")
	p.NewBody()

	kubecontext := "blacontext"
	version := 666
	errExpected := errors.New("test error")

	mockedRelease := NewMockReleaseConfig(s.T())
	mockedRelease.On("KubeContext").Return(kubecontext)
	mockedRelease.On("Get", version).Return((*helmRelease.Release)(nil), errExpected)
	p.SetReleases(mockedRelease)

	_, _, err := p.trackerRollbackObjects(version, &tracker.Config{TrackGeneric: true})

	s.Require().ErrorIs(err, errExpected)
	mockedRelease.AssertExpectations(s.T())
}

func (s *TrackerTestSuite) TestMultipleContexts() {
	p := New("")
	p.NewBody()

	relName := "bla"
	relNS := "blabla"
	u, _ := uniqname.New(relName, relNS, "")

	mockedRelease1 := NewMockReleaseConfig(s.T())
	mockedRelease1.On("KubeContext").Return("bla1")
	mockedRelease1.On("Uniq").Return(u)
	mockedRelease1.On("Namespace").Return(relNS)

	mockedRelease2 := NewMockReleaseConfig(s.T())
	mockedRelease2.On("KubeContext").Return("bla2")
	mockedRelease2.On("Uniq").Return(u)
	mockedRelease2.On("Namespace").Return(relNS)

	p.SetReleases(mockedRelease1, mockedRelease2)

	_, _, err := p.trackerObjects(&tracker.Config{TrackGeneric: true}, func(rel release.Config) (string, error) {
		return "", nil
	})

	s.Require().ErrorIs(err, ErrMultipleKubecontexts)
	mockedRelease1.AssertExpectations(s.T())
	mockedRelease2.AssertExpectations(s.T())
}
