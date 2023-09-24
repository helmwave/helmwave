package plan

import (
	"fmt"
	"testing"

	"github.com/helmwave/helmwave/pkg/registry"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	helmRegistry "helm.sh/helm/v3/pkg/registry"
)

type BuildRegistriesTestSuite struct {
	suite.Suite
}

func TestBuildRegistriesTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildRegistriesTestSuite))
}

func (ts *BuildRegistriesTestSuite) TestUnusedRegistry() {
	p := New()

	regi := &MockRegistryConfig{}
	p.SetRegistries(regi)

	regis, err := p.buildRegistries()
	ts.Require().NoError(err)
	ts.Require().Empty(regis)

	regi.AssertExpectations(ts.T())
}

func (ts *BuildRegistriesTestSuite) TestNoOCIRegistries() {
	p := New()

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Chart").Return(&release.Chart{})

	p.SetReleases(mockedRelease)

	repos, err := p.buildRegistries()
	ts.Require().NoError(err)
	ts.Require().Empty(repos)

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *BuildRegistriesTestSuite) TestMissingRegistry() {
	p := New()

	regiName := "blablanami"

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Repo").Return(regiName)
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Logger").Return(log.WithField("test", ts.T().Name()))
	mockedRelease.On("Chart").Return(&release.Chart{Name: fmt.Sprintf("%s://%s", helmRegistry.OCIScheme, regiName)})

	p.SetReleases(mockedRelease)

	repos, err := p.buildRegistries()

	ts.Require().ErrorIs(err, registry.NotFoundError{})
	ts.Require().Empty(repos)

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *BuildRegistriesTestSuite) TestSuccess() {
	p := New()

	regiHost := "blablanami"

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Repo").Return(regiHost)
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Logger").Return(log.WithField("test", ts.T().Name()))
	mockedRelease.On("Chart").Return(&release.Chart{Name: fmt.Sprintf("%s://", helmRegistry.OCIScheme)})

	regi := &MockRegistryConfig{}
	regi.On("Host").Return(regiHost)

	p.SetReleases(mockedRelease)
	p.SetRegistries(regi)

	repos, err := p.buildRegistries()

	ts.Require().NoError(err)
	ts.Require().Len(repos, 1)
	ts.Require().Contains(repos, regi)

	mockedRelease.AssertExpectations(ts.T())
	regi.AssertExpectations(ts.T())
}
