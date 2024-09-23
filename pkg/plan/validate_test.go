package plan_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/fileref"

	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/tests"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/stretchr/testify/suite"
)

type ValidateTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestValidateTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValidateTestSuite))
}

func (ts *ValidateTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *ValidateTestSuite) TestInvalidRelease() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	err := errors.New("test error")

	mockedRelease := plan.NewMockReleaseConfig(ts.T())
	mockedRelease.On("Validate").Return(err)

	p.SetReleases(mockedRelease)

	ts.Require().ErrorIs(err, body.ValidateReleases())
	ts.Require().ErrorIs(err, body.Validate())

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *ValidateTestSuite) TestInvalidRepository() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	err := errors.New("test error")

	mockedRepo := plan.NewMockRepositoryConfig(ts.T())
	mockedRepo.On("Validate").Return(err)

	p.SetRepositories(mockedRepo)

	ts.Require().ErrorIs(err, body.ValidateRepositories())
	ts.Require().ErrorIs(err, body.Validate())

	mockedRepo.AssertExpectations(ts.T())
}

func (ts *ValidateTestSuite) TestValidateValues() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	valuesContents := []byte("a: b")
	tmpValues := filepath.Join(tmpDir, "valuesName")
	ts.Require().NoError(os.WriteFile(tmpValues, valuesContents, 0o600))

	mockedRelease := plan.NewMockReleaseConfig(ts.T())
	mockedRelease.On("Name").Return(ts.T().Name())
	mockedRelease.On("Namespace").Return(ts.T().Name())
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("KubeContext").Return("")

	v := fileref.Config{Src: tmpValues}
	ts.Require().NoError(v.Set(context.Background(), "test.values.yml", tmpDir, template.TemplaterSprig, nil))

	mockedRelease.On("Values").Return([]fileref.Config{v})

	p.SetReleases(mockedRelease)

	ts.Require().NoError(p.ValidateValuesImport())

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *ValidateTestSuite) TestValidateValuesNotFound() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	valuesContents := []byte("a: b")
	tmpValues := filepath.Join(tmpDir, "valuesName")
	ts.Require().NoError(os.WriteFile(tmpValues, valuesContents, 0o600))

	mockedRelease := plan.NewMockReleaseConfig(ts.T())
	v := fileref.Config{Src: tmpValues}
	mockedRelease.On("Values").Return([]fileref.Config{v})

	p.SetReleases(mockedRelease)

	ts.Require().Error(p.ValidateValuesImport())

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *ValidateTestSuite) TestValidateValuesNoReleases() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	p.NewBody()

	ts.Require().NoError(p.ValidateValuesImport())
}

func (ts *ValidateTestSuite) TestValidateRepositoryDuplicate() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	mockedRepo := plan.NewMockRepositoryConfig(ts.T())
	mockedRepo.On("Name").Return("blabla")
	mockedRepo.On("Validate").Return(nil)

	p.SetRepositories(mockedRepo, mockedRepo)

	var e *repo.DuplicateError

	ts.Require().ErrorAs(body.ValidateRepositories(), &e)
	ts.Equal("blabla", e.Name)

	ts.Require().ErrorAs(body.Validate(), &e)
	ts.Equal("blabla", e.Name)

	mockedRepo.AssertExpectations(ts.T())
}

func (ts *ValidateTestSuite) TestValidateReleaseDuplicate() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	mockedRelease := plan.NewMockReleaseConfig(ts.T())
	mockedRelease.On("Name").Return("blabla")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Validate").Return(nil)
	mockedRelease.On("KubeContext").Return("")

	p.SetReleases(mockedRelease, mockedRelease)

	var e *release.DuplicateError

	ts.Require().ErrorAs(body.ValidateReleases(), &e)
	ts.Equal(mockedRelease.Uniq(), e.Uniq)

	ts.Require().ErrorAs(body.Validate(), &e)
	ts.Equal(mockedRelease.Uniq(), e.Uniq)

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *ValidateTestSuite) TestValidateEmpty() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	body := p.NewBody()

	ts.Require().NoError(body.Validate())
}
