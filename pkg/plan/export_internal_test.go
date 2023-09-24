package plan

import (
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type ExportTestSuite struct {
	suite.Suite
}

func (s *ExportTestSuite) TestValuesEmpty() {
	p := New()
	p.templater = template.TemplaterSprig

	p.body = &planBody{}

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := p.exportValues(baseFS.(fs.StatFS), baseFS.(fsimpl.WriteableFS))
	s.Require().NoError(err)
}

func (s *ExportTestSuite) TestValuesOneRelease() {
	tmpDir := s.T().TempDir()
	p := New()
	p.templater = template.TemplaterSprig

	baseFS1, _ := filefs.New(&url.URL{Scheme: "file", Path: tmpDir})
	baseFS2, _ := filefs.New(&url.URL{Scheme: "file", Path: s.T().TempDir()})

	valuesName := "blablavalues.yaml"
	valuesContents := []byte("a: b")

	f, err := helper.CreateFile(baseFS1.(fsimpl.WriteableFS), valuesName)
	s.Require().NoError(err)
	_, err = f.Write(valuesContents)
	s.Require().NoError(err)
	s.Require().NoError(f.Close())

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Values").Return([]release.ValuesReference{
		{Src: valuesName},
	})
	mockedRelease.On("ExportValues").Return(nil)
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))

	p.body = &planBody{
		Releases: release.Configs{mockedRelease},
	}

	s.Require().NoError(p.exportValues(baseFS1.(fs.StatFS), baseFS2.(fsimpl.WriteableFS)))
	mockedRelease.AssertExpectations(s.T())
	s.Require().FileExists(filepath.Join(baseFS2.(fsimpl.CurrentPathFS).CurrentPath(), Values, valuesName))

	contents, err := os.ReadFile(filepath.Join(baseFS2.(fsimpl.CurrentPathFS).CurrentPath(), Values, valuesName))
	s.Require().NoError(err)
	s.Require().Equal(valuesContents, contents)
}

func TestExportTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExportTestSuite))
}
