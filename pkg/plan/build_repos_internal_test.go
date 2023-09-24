package plan

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
)

type BuildRepositoriesTestSuite struct {
	suite.Suite
}

func TestBuildRepositoriesTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildRepositoriesTestSuite))
}

func (ts *BuildRepositoriesTestSuite) TestReposEmpty() {
	p := New()

	p.NewBody()

	repos, err := p.buildRepositories()
	ts.Require().NoError(err)
	ts.Require().Empty(repos)
}

func (ts *BuildRepositoriesTestSuite) TestLocalRepo() {
	p := New()

	repoName := ""

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Repo").Return(repoName)
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Chart").Return(&release.Chart{})

	mockedRepo := &MockRepositoryConfig{}

	p.SetRepositories(mockedRepo)
	p.SetReleases(mockedRelease)

	repos, err := p.buildRepositories()
	ts.Require().NoError(err)
	ts.Require().Empty(repos, 0)

	mockedRepo.AssertExpectations(ts.T())
	mockedRelease.AssertExpectations(ts.T())
}

func (ts *BuildRepositoriesTestSuite) TestUnusedRepo() {
	p := New()

	mockedRepo := &MockRepositoryConfig{}

	p.SetRepositories(mockedRepo)

	repos, err := p.buildRepositories()
	ts.Require().NoError(err)
	ts.Require().Empty(repos)

	mockedRepo.AssertExpectations(ts.T())
}

func (ts *BuildRepositoriesTestSuite) TestSuccess() {
	p := New()

	repoName := "blablanami"

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Repo").Return(repoName)
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Chart").Return(&release.Chart{})

	mockedRepo := &MockRepositoryConfig{}
	mockedRepo.On("Name").Return(repoName)

	p.SetRepositories(mockedRepo)
	p.SetReleases(mockedRelease)

	repos, err := p.buildRepositories()
	ts.Require().NoError(err)
	ts.Require().Len(repos, 1)
	ts.Require().Contains(repos, mockedRepo)

	mockedRepo.AssertExpectations(ts.T())
	mockedRelease.AssertExpectations(ts.T())
}

func (ts *BuildRepositoriesTestSuite) TestMissingRepo() {
	p := New()

	repoName := "blablanami"

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Repo").Return(repoName)
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Chart").Return(&release.Chart{})

	p.SetReleases(mockedRelease)

	repos, err := p.buildRepositories()
	ts.Require().Error(err)
	ts.Require().Empty(repos)

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *BuildRepositoriesTestSuite) TestRepoIsLocal() {
	ts.Require().True(repoIsLocal(ts.T().TempDir()))
}
