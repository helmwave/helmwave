package plan

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
)

type ExportTestSuite struct {
	suite.Suite
}

func (s *ExportTestSuite) TestValuesEmpty() {
	tmpDir := s.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))
	p.templater = "sprig"

	p.body = &planBody{}

	err := p.exportValues()
	s.Require().NoError(err)
}

func (s *ExportTestSuite) TestValuesOneRelease() {
	tmpDir := s.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))
	p.templater = "sprig"

	valuesName := "blablavalues.yaml"
	valuesContents := []byte("a: b")
	tmpValues := filepath.Join(tmpDir, valuesName)
	s.Require().NoError(os.WriteFile(tmpValues, valuesContents, 0o600))

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Values").Return([]release.ValuesReference{
		{Src: tmpValues},
	})
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("BuildValues").Return(nil)
	mockedRelease.On("Uniq").Return()

	p.body = &planBody{
		Releases: release.Configs{mockedRelease},
	}

	s.Require().NoError(p.buildValues())
	s.Require().NoError(p.exportValues())
	mockedRelease.AssertExpectations(s.T())
	s.Require().DirExists(filepath.Join(tmpDir, Dir, Values))
	s.Require().FileExists(filepath.Join(tmpDir, Dir, Values, valuesName))

	contents, err := os.ReadFile(filepath.Join(tmpDir, Dir, Values, valuesName))
	s.Require().NoError(err)
	s.Require().Equal(valuesContents, contents)
}

func TestExportTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExportTestSuite))
}
