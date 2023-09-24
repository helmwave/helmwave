package plan

import (
	"errors"
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

type BuildValuesTestSuite struct {
	suite.Suite
}

func (s *BuildValuesTestSuite) createPlan() *Plan {
	s.T().Helper()

	p := New()
	p.templater = template.TemplaterSprig

	return p
}

func (s *BuildValuesTestSuite) TestValuesEmpty() {
	p := s.createPlan()

	p.body = &planBody{}

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	s.Require().NoError(p.buildValues(baseFS.(fs.StatFS), baseFS.(fsimpl.WriteableFS)))
}

func (s *BuildValuesTestSuite) TestValuesBuildError() {
	tmpDir := s.T().TempDir()

	p := s.createPlan()

	tmpValues := filepath.Join(tmpDir, "blablavalues.yaml")
	s.Require().NoError(os.WriteFile(tmpValues, []byte("a: b"), 0o600))

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Values").Return([]release.ValuesReference{
		{Src: tmpValues},
	})

	errBuildValues := errors.New("values build error")
	mockedRelease.On("ExportValues").Return(errBuildValues)

	p.body = &planBody{
		Releases: release.Configs{mockedRelease},
	}

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	s.Require().ErrorIs(p.buildValues(baseFS.(fs.StatFS), baseFS.(fsimpl.WriteableFS)), errBuildValues)
	mockedRelease.AssertExpectations(s.T())
}

func (s *BuildValuesTestSuite) TestSuccess() {
	tmpDir := s.T().TempDir()
	p := s.createPlan()

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
	mockedRelease.On("ExportValues").Return(nil)
	mockedRelease.On("Uniq").Return()

	p.body = &planBody{
		Releases: release.Configs{mockedRelease},
	}

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	s.Require().NoError(p.buildValues(baseFS.(fs.StatFS), baseFS.(fsimpl.WriteableFS)))
	mockedRelease.AssertExpectations(s.T())
}

func TestBuildValuesTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildValuesTestSuite))
}
