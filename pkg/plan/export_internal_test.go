package plan

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/fileref"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type ExportTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestExportTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExportTestSuite))
}

func (ts *ExportTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *ExportTestSuite) TestValuesEmpty() {
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))
	p.templater = template.TemplaterSprig

	p.body = &planBody{}

	err := p.exportValues()
	ts.Require().NoError(err)
}

func (ts *ExportTestSuite) TestValuesOneRelease() {
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))
	p.templater = template.TemplaterSprig

	valuesName := "blablavalues.yaml"
	valuesContents := []byte("a: b")
	tmpValues := filepath.Join(tmpDir, valuesName)
	ts.Require().NoError(os.WriteFile(tmpValues, valuesContents, 0o600))

	mockedRelease := NewMockReleaseConfig(ts.T())
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Values").Return([]fileref.Config{
		{Src: tmpValues},
	})
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("BuildValues").Return(nil)
	mockedRelease.On("KubeContext").Return("")
	mockedRelease.On("Uniq").Return()

	p.body = &planBody{
		Releases: release.Configs{mockedRelease},
	}

	ts.Require().NoError(p.buildValues(ts.ctx))
	ts.Require().NoError(p.exportValues())
	mockedRelease.AssertExpectations(ts.T())
	ts.Require().DirExists(filepath.Join(tmpDir, Dir, Values))
	ts.Require().FileExists(filepath.Join(tmpDir, Dir, Values, valuesName))

	contents, err := os.ReadFile(filepath.Join(tmpDir, Dir, Values, valuesName))
	ts.Require().NoError(err)
	ts.Require().Equal(valuesContents, contents)
}
