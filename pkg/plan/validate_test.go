package plan_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type ValidateTestSuite struct {
	suite.Suite
}

func (s *ValidateTestSuite) TestValidateValues() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	valuesContents := []byte("a: b")
	tmpValues := filepath.Join(tmpDir, "valuesName")
	s.Require().NoError(os.WriteFile(tmpValues, valuesContents, 0o600))

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Name").Return(s.T().Name())
	mockedRelease.On("Namespace").Return(s.T().Name())
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))
	v := release.ValuesReference{Src: tmpValues}
	s.Require().NoError(v.SetViaRelease(mockedRelease, tmpDir, "sprig"))
	mockedRelease.On("Values").Return([]release.ValuesReference{v})

	p.SetReleases(mockedRelease)

	s.Require().NoError(p.ValidateValuesImport())

	mockedRelease.AssertExpectations(s.T())
}

func (s *ValidateTestSuite) TestValidateValuesNotFound() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	valuesContents := []byte("a: b")
	tmpValues := filepath.Join(tmpDir, "valuesName")
	s.Require().NoError(os.WriteFile(tmpValues, valuesContents, 0o600))

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))
	v := release.ValuesReference{Src: tmpValues}
	mockedRelease.On("Values").Return([]release.ValuesReference{v})

	p.SetReleases(mockedRelease)

	s.Require().Error(p.ValidateValuesImport())

	mockedRelease.AssertExpectations(s.T())
}

func (s *ValidateTestSuite) TestValidateValuesNoReleases() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	p.NewBody()

	s.Require().NoError(p.ValidateValuesImport())
}

func (s *ValidateTestSuite) TestValidateRepositoryName() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	mockedRepo := &plan.MockRepoConfig{}
	mockedRepo.On("Name").Return("")

	p.SetRepositories(mockedRepo)

	s.Require().Error(body.ValidateRepositories())
	s.Require().Error(body.Validate())

	mockedRepo.AssertExpectations(s.T())
}

func (s *ValidateTestSuite) TestValidateRepositoryURLEmpty() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	mockedRepo := &plan.MockRepoConfig{}
	mockedRepo.On("Name").Return("blabla")
	mockedRepo.On("URL").Return("")

	p.SetRepositories(mockedRepo)

	s.Require().Error(body.ValidateRepositories())
	s.Require().Error(body.Validate())

	mockedRepo.AssertExpectations(s.T())
}

func (s *ValidateTestSuite) TestValidateRepositoryURL() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	mockedRepo := &plan.MockRepoConfig{}
	mockedRepo.On("Name").Return("blabla")
	mockedRepo.On("URL").Return("\\asdasd://null")

	p.SetRepositories(mockedRepo)

	s.Require().Error(body.ValidateRepositories())
	s.Require().Error(body.Validate())

	mockedRepo.AssertExpectations(s.T())
}

func (s *ValidateTestSuite) TestValidateRepositoryDuplicate() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	mockedRepo := &plan.MockRepoConfig{}
	mockedRepo.On("Name").Return("blabla")
	mockedRepo.On("URL").Return("http://localhost")

	p.SetRepositories(mockedRepo, mockedRepo)

	s.Require().Error(body.ValidateRepositories())
	s.Require().Error(body.Validate())

	mockedRepo.AssertExpectations(s.T())
}

func (s *ValidateTestSuite) TestValidateRepository() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	mockedRepo := &plan.MockRepoConfig{}
	mockedRepo.On("Name").Return("blabla")
	mockedRepo.On("URL").Return("http://localhost")

	p.SetRepositories(mockedRepo)

	s.Require().NoError(body.ValidateRepositories())
	s.Require().NoError(body.Validate())

	mockedRepo.AssertExpectations(s.T())
}

func (s *ValidateTestSuite) TestValidateReleaseName() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Name").Return("")

	p.SetReleases(mockedRelease)

	s.Require().Error(body.ValidateReleases())
	s.Require().Error(body.Validate())

	mockedRelease.AssertExpectations(s.T())
}

func (s *ValidateTestSuite) TestValidateReleaseNamespace() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Name").Return("blabla")
	mockedRelease.On("Namespace").Return("///")

	p.SetReleases(mockedRelease)

	s.Require().Error(body.ValidateReleases())
	s.Require().Error(body.Validate())

	mockedRelease.AssertExpectations(s.T())
}

func (s *ValidateTestSuite) TestValidateReleaseUniq() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Name").Return("blabla")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()

	p.SetReleases(mockedRelease)

	s.Require().NoError(body.ValidateReleases())
	s.Require().NoError(body.Validate())

	mockedRelease.AssertExpectations(s.T())
}

func (s *ValidateTestSuite) TestValidateReleaseDuplicate() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Name").Return("blabla")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()

	p.SetReleases(mockedRelease, mockedRelease)

	s.Require().Error(body.ValidateReleases())
	s.Require().Error(body.Validate())

	mockedRelease.AssertExpectations(s.T())
}

func (s *ValidateTestSuite) TestValidateEmpty() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	err := body.Validate()
	s.Require().Error(err)
}

func TestValidateTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValidateTestSuite))
}
