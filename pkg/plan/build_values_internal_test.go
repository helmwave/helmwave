package plan

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/fileref"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type BuildValuesTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestBuildValuesTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildValuesTestSuite))
}

func (ts *BuildValuesTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *BuildValuesTestSuite) createPlan(tmpDir string) *Plan {
	ts.T().Helper()

	p := New(filepath.Join(tmpDir, Dir))
	p.templater = template.TemplaterSprig

	return p
}

func (ts *BuildValuesTestSuite) TestValuesEmpty() {
	tmpDir := ts.T().TempDir()
	p := ts.createPlan(tmpDir)

	p.body = &planBody{}

	ts.Require().NoError(p.buildValues(ts.ctx))
}

func (ts *BuildValuesTestSuite) TestValuesBuildError() {
	tmpDir := ts.T().TempDir()
	p := ts.createPlan(tmpDir)

	tmpValues := filepath.Join(tmpDir, "blablavalues.yaml")
	ts.Require().NoError(os.WriteFile(tmpValues, []byte("a: b"), 0o600))

	mockedRelease := NewMockReleaseConfig(ts.T())
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("KubeContext").Return("")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Values").Return([]fileref.Config{
		{Src: tmpValues},
	})

	errBuildValues := errors.New("values build error")
	mockedRelease.On("BuildValues").Return(errBuildValues)

	p.body = &planBody{
		Releases: release.Configs{mockedRelease},
	}

	ts.Require().ErrorIs(p.buildValues(ts.ctx), errBuildValues)
	mockedRelease.AssertExpectations(ts.T())
}

func (ts *BuildValuesTestSuite) TestSuccess() {
	tmpDir := ts.T().TempDir()
	p := ts.createPlan(tmpDir)

	valuesName := "blablavalues.yaml"
	valuesContents := []byte("a: b")
	tmpValues := filepath.Join(tmpDir, valuesName)
	ts.Require().NoError(os.WriteFile(tmpValues, valuesContents, 0o600))

	mockedRelease := NewMockReleaseConfig(ts.T())
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("KubeContext").Return("")
	mockedRelease.On("Values").Return([]fileref.Config{
		{Src: tmpValues},
	})
	mockedRelease.On("BuildValues").Return(nil)
	mockedRelease.On("Uniq").Return()

	p.body = &planBody{
		Releases: release.Configs{mockedRelease},
	}

	ts.Require().NoError(p.buildValues(ts.ctx))
	mockedRelease.AssertExpectations(ts.T())
}
