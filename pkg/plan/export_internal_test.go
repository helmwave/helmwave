package plan

import (
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/stretchr/testify/suite"
)

type ExportTestSuite struct {
	suite.Suite
}

func (s *ExportTestSuite) TestValuesEmpty() {
	tmpDir := s.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))
	p.templater = template.TemplaterSprig

	p.body = &planBody{}

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := p.exportValues(baseFS.(fsimpl.WriteableFS))
	s.Require().NoError(err)
}

func (s *ExportTestSuite) TestValuesOneRelease() {
	tmpDir := s.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))
	p.templater = template.TemplaterSprig

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

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	s.Require().NoError(p.buildValues(baseFS.(fs.StatFS), baseFS.(fsimpl.WriteableFS)))
	s.Require().NoError(p.exportValues(baseFS.(fsimpl.WriteableFS)))
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
