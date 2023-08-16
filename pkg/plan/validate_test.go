package plan_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/repo"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type ValidateTestSuite struct {
	suite.Suite
}

func (s *ValidateTestSuite) TestInvalidRelease() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	err := errors.New("test error")

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Validate").Return(err)

	p.SetReleases(mockedRelease)

	s.Require().ErrorIs(err, body.ValidateReleases())
	s.Require().Error(err, body.Validate())

	mockedRelease.AssertExpectations(s.T())
}

func (s *ValidateTestSuite) TestInvalidRepository() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	err := errors.New("test error")

	mockedRepo := &plan.MockRepoConfig{}
	mockedRepo.On("Validate").Return(err)

	p.SetRepositories(mockedRepo)

	s.Require().ErrorIs(err, body.ValidateRepositories())
	s.Require().Error(err, body.Validate())

	mockedRepo.AssertExpectations(s.T())
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
	s.Require().NoError(v.SetViaRelease(mockedRelease, tmpDir, template.TemplaterSprig))
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

func (s *ValidateTestSuite) TestValidateRepositoryDuplicate() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	mockedRepo := &plan.MockRepoConfig{}
	mockedRepo.On("Name").Return("blabla")
	mockedRepo.On("Validate").Return(nil)

	p.SetRepositories(mockedRepo, mockedRepo)

	s.Require().ErrorIs(repo.DuplicateError{}, body.ValidateRepositories())
	s.Require().ErrorIs(repo.DuplicateError{}, body.Validate())

	mockedRepo.AssertExpectations(s.T())
}

func (s *ValidateTestSuite) TestValidateReleaseDuplicate() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Name").Return("blabla")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Validate").Return(nil)

	p.SetReleases(mockedRelease, mockedRelease)

	s.Require().ErrorIs(release.DuplicateError{}, body.ValidateReleases())
	s.Require().ErrorIs(release.DuplicateError{}, body.Validate())

	mockedRelease.AssertExpectations(s.T())
}

func (s *ValidateTestSuite) TestValidateEmpty() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	err := body.Validate()
	s.Require().ErrorIs(plan.ErrEmptyPlan, err)
}

func TestValidateTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValidateTestSuite))
}
