package plan

import (
	"errors"
	"testing"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

type KubedogTestSuite struct {
	suite.Suite
}

func TestKubedogTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(KubedogTestSuite))
}

func (s *KubedogTestSuite) TestNoReleases() {
	p := New("")
	p.NewBody()

	spec, _, err := p.kubedogSpecs(&kubedog.Config{}, nil)

	s.Require().NoError(err)

	s.Require().Empty(spec.Generics)
	s.Require().Empty(spec.Jobs)
	s.Require().Empty(spec.DaemonSets)
	s.Require().Empty(spec.Canaries)
	s.Require().Empty(spec.Deployments)
	s.Require().Empty(spec.StatefulSets)
}

func (s *KubedogTestSuite) TestCallsManifestFunction() {
	p := New("")
	p.NewBody()

	rel := NewMockReleaseConfig(s.T())
	p.SetReleases(rel)

	s.Require().Panics(func() {
		_, _, _ = p.kubedogSpecs(&kubedog.Config{}, nil)
	})
}

func (s *KubedogTestSuite) TestSyncSpecs() {
	p := New("")
	p.NewBody()

	relName := "bla"
	relNS := "blabla"
	kubecontext := "blacontext"
	u, _ := uniqname.New(relName, relNS, "")

	p.manifests[u] = `
kind: Canary
metadata:
  name: blabla
---
kind: DaemonSet
metadata:
  name: blabla
---
kind: Deployment
metadata:
  name: blabla
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: bla
---
kind: Job
metadata:
  name: blabla
---
kind: StatefulSet
metadata:
  name: blabla
---
`

	mockedRelease := NewMockReleaseConfig(s.T())
	mockedRelease.On("KubeContext").Return(kubecontext)
	mockedRelease.On("Uniq").Return(u)
	mockedRelease.On("Namespace").Return(relNS)
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))
	p.SetReleases(mockedRelease)

	spec, context, err := p.kubedogSyncSpecs(&kubedog.Config{TrackGeneric: true})

	s.Require().NoError(err)

	s.Require().Len(spec.Canaries, 1)
	s.Require().Len(spec.DaemonSets, 1)
	s.Require().Len(spec.Deployments, 1)
	s.Require().Len(spec.Generics, 1)
	s.Require().Len(spec.Jobs, 1)
	s.Require().Len(spec.StatefulSets, 1)

	s.Require().Equal(kubecontext, context)

	mockedRelease.AssertExpectations(s.T())
}

func (s *KubedogTestSuite) TestRollbackSpecs() {
	p := New("")
	p.NewBody()

	relName := "bla"
	relNS := "blabla"
	kubecontext := "blacontext"
	version := 666
	u, _ := uniqname.New(relName, relNS, "")

	manifest := `
kind: Canary
metadata:
  name: blabla
---
kind: DaemonSet
metadata:
  name: blabla
---
kind: Deployment
metadata:
  name: blabla
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: bla
---
kind: Job
metadata:
  name: blabla
---
kind: StatefulSet
metadata:
  name: blabla
---
`

	mockedRelease := NewMockReleaseConfig(s.T())
	mockedRelease.On("KubeContext").Return(kubecontext)
	mockedRelease.On("Uniq").Return(u)
	mockedRelease.On("Namespace").Return(relNS)
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))
	mockedRelease.On("Get", version).Return(&helmRelease.Release{Manifest: manifest}, nil)
	p.SetReleases(mockedRelease)

	spec, context, err := p.kubedogRollbackSpecs(version, &kubedog.Config{TrackGeneric: true})

	s.Require().NoError(err)

	s.Require().Len(spec.Canaries, 1)
	s.Require().Len(spec.DaemonSets, 1)
	s.Require().Len(spec.Deployments, 1)
	s.Require().Len(spec.Generics, 1)
	s.Require().Len(spec.Jobs, 1)
	s.Require().Len(spec.StatefulSets, 1)

	s.Require().Equal(kubecontext, context)

	mockedRelease.AssertExpectations(s.T())
}

func (s *KubedogTestSuite) TestRollbackSpecsGetError() {
	p := New("")
	p.NewBody()

	kubecontext := "blacontext"
	version := 666
	errExpected := errors.New("test error")

	mockedRelease := NewMockReleaseConfig(s.T())
	mockedRelease.On("KubeContext").Return(kubecontext)
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))
	mockedRelease.On("Get", version).Return((*helmRelease.Release)(nil), errExpected)
	p.SetReleases(mockedRelease)

	_, _, err := p.kubedogRollbackSpecs(version, &kubedog.Config{TrackGeneric: true})

	s.Require().ErrorIs(err, errExpected)
	mockedRelease.AssertExpectations(s.T())
}

func (s *KubedogTestSuite) TestSpecsMultipleContexts() {
	p := New("")
	p.NewBody()

	relName := "bla"
	relNS := "blabla"
	u, _ := uniqname.New(relName, relNS, "")

	mockedRelease1 := NewMockReleaseConfig(s.T())
	mockedRelease1.On("KubeContext").Return("bla1")
	mockedRelease1.On("Uniq").Return(u)
	mockedRelease1.On("Namespace").Return(relNS)
	mockedRelease1.On("Logger").Return(log.WithField("test", s.T().Name()))

	mockedRelease2 := NewMockReleaseConfig(s.T())
	mockedRelease2.On("KubeContext").Return("bla2")
	mockedRelease2.On("Uniq").Return(u)
	mockedRelease2.On("Namespace").Return(relNS)
	mockedRelease2.On("Logger").Return(log.WithField("test", s.T().Name()))

	p.SetReleases(mockedRelease1, mockedRelease2)

	_, _, err := p.kubedogSpecs(&kubedog.Config{TrackGeneric: true}, func(rel release.Config) (string, error) {
		return "", nil
	})

	s.Require().ErrorIs(err, ErrMultipleKubecontexts)
	mockedRelease1.AssertExpectations(s.T())
	mockedRelease2.AssertExpectations(s.T())
}
