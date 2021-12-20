package plan

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BuildRepositoriesTestSuite struct {
	suite.Suite
}

func (s *BuildRepositoriesTestSuite) TestReposEmpty() {
	tmpDir := s.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	p.body = &planBody{}

	repos, err := p.buildRepositories()
	s.Require().NoError(err)
	s.Require().Empty(repos)
}

func (s *BuildRepositoriesTestSuite) TestUnusedRepo() {
	tmpDir := s.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	mockedRepo := &mockRepoConfig{}

	p.body = &planBody{
		Repositories: repoConfigs{mockedRepo},
	}

	repos, err := p.buildRepositories()
	s.Require().NoError(err)
	s.Require().Empty(repos)

	mockedRepo.AssertExpectations(s.T())
}

func (s *BuildRepositoriesTestSuite) TestSuccess() {
	tmpDir := s.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	repoName := "blablanami"

	mockedRelease := &mockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Repo").Return(repoName)
	mockedRelease.On("Namespace").Return("defaultblabla")

	mockedRepo := &mockRepoConfig{}
	mockedRepo.On("Name").Return(repoName)

	p.body = &planBody{
		Repositories: repoConfigs{mockedRepo},
		Releases:     releaseConfigs{mockedRelease},
	}

	repos, err := p.buildRepositories()
	s.Require().NoError(err)
	s.Require().Len(repos, 1)
	s.Require().Contains(repos, mockedRepo)

	mockedRepo.AssertExpectations(s.T())
	mockedRelease.AssertExpectations(s.T())
}

func (s *BuildRepositoriesTestSuite) TestMissingRepo() {
	tmpDir := s.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	repoName := "blablanami"

	mockedRelease := &mockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Repo").Return(repoName)
	mockedRelease.On("Namespace").Return("defaultblabla")

	p.body = &planBody{
		Releases: releaseConfigs{mockedRelease},
	}

	repos, err := p.buildRepositories()
	s.Require().Error(err)
	s.Require().Empty(repos)

	mockedRelease.AssertExpectations(s.T())
}

func TestBuildRepositoriesTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildRepositoriesTestSuite))
}
