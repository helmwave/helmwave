package plan

import (
	"path/filepath"
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
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	p.NewBody()

	repos, err := p.buildRepositories()
	ts.Require().NoError(err)
	ts.Require().Empty(repos)
}

func (ts *BuildRepositoriesTestSuite) TestLocalRepo() {
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	repoName := ts.T().Name()

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Repo").Return(repoName)
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Chart").Return(&release.Chart{})
	mockedRelease.On("KubeContext").Return("")

	mockedRepo := &MockRepositoryConfig{}
	mockedRepo.On("Name").Return(repoName)

	p.SetRepositories(mockedRepo)
	p.SetReleases(mockedRelease)

	repos, err := p.buildRepositories()
	ts.Require().NoError(err)
	ts.Len(repos, 1)
	ts.Contains(repos, mockedRepo)

	mockedRepo.AssertExpectations(ts.T())
	mockedRelease.AssertExpectations(ts.T())
}

func (ts *BuildRepositoriesTestSuite) TestUnusedRepo() {
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	mockedRepo := &MockRepositoryConfig{}

	p.SetRepositories(mockedRepo)

	repos, err := p.buildRepositories()
	ts.Require().NoError(err)
	ts.Empty(repos)

	mockedRepo.AssertExpectations(ts.T())
}

func (ts *BuildRepositoriesTestSuite) TestSuccess() {
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	repoName := "blablanami"

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Repo").Return(repoName)
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Chart").Return(&release.Chart{})
	mockedRelease.On("KubeContext").Return("")

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
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	repoName := "blablanami"

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Repo").Return(repoName)
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Chart").Return(&release.Chart{})
	mockedRelease.On("KubeContext").Return("")

	p.SetReleases(mockedRelease)

	repos, err := p.buildRepositories()
	ts.Require().Error(err)
	ts.Require().Empty(repos)

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *BuildRepositoriesTestSuite) TestRepoIsLocal() {
	ts.Require().True(repoIsLocal(ts.T().TempDir()))
}
