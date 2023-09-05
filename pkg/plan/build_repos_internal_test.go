package plan

import (
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/stretchr/testify/suite"
)

type BuildRepositoriesTestSuite struct {
	suite.Suite
}

func TestBuildRepositoriesTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildRepositoriesTestSuite))
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

	mockedRepo := &MockRepoConfig{}

	p.body = &planBody{
		Repositories: repo.Configs{mockedRepo},
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

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Repo").Return(repoName)
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Chart").Return(&release.Chart{})

	mockedRepo := &MockRepoConfig{}
	mockedRepo.On("Name").Return(repoName)

	p.body = &planBody{
		Repositories: repo.Configs{mockedRepo},
		Releases:     release.Configs{mockedRelease},
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

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Repo").Return(repoName)
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Chart").Return(&release.Chart{})

	p.body = &planBody{
		Releases: release.Configs{mockedRelease},
	}

	repos, err := p.buildRepositories()
	s.Require().Error(err)
	s.Require().Empty(repos)

	mockedRelease.AssertExpectations(s.T())
}
