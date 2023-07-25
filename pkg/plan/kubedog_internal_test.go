package plan

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v3/pkg/release"
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

	rel := &MockReleaseConfig{}
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
	u, _ := uniqname.Generate(relName, relNS)

	p.manifests[u] = `
kind: Canary
---
kind: DaemonSet
---
kind: Deployment
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: bla
---
kind: Job
---
kind: StatefulSet
---
`

	mockedRelease := &MockReleaseConfig{}
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
	u, _ := uniqname.Generate(relName, relNS)

	manifest := `
kind: Canary
---
kind: DaemonSet
---
kind: Deployment
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: bla
---
kind: Job
---
kind: StatefulSet
---
`

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("KubeContext").Return(kubecontext)
	mockedRelease.On("Uniq").Return(u)
	mockedRelease.On("Namespace").Return(relNS)
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))
	mockedRelease.On("Get", version).Return(&release.Release{Manifest: manifest}, nil)
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
