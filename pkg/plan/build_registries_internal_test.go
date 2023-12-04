package plan

import (
	"fmt"
	"path/filepath"
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
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	regi := &MockRegistryConfig{}
	p.SetRegistries(regi)

	regis, err := p.buildRegistries()
	ts.Require().NoError(err)
	ts.Require().Empty(regis)

	regi.AssertExpectations(ts.T())
}

func (ts *BuildRegistriesTestSuite) TestNoOCIRegistries() {
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Chart").Return(&release.Chart{})

	p.SetReleases(mockedRelease)

	repos, err := p.buildRegistries()
	ts.Require().NoError(err)
	ts.Require().Empty(repos)

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *BuildRegistriesTestSuite) TestMissingRegistry() {
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

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

	var e *registry.NotFoundError
	ts.Require().ErrorAs(err, &e)
	ts.Equal(regiName, e.Host)

	ts.Empty(repos)

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *BuildRegistriesTestSuite) TestSuccess() {
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

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
