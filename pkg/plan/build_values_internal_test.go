package plan

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
)

type BuildValuesTestSuite struct {
	suite.Suite
}

func (s *BuildValuesTestSuite) TestValuesEmpty() {
	tmpDir := s.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	p.body = &planBody{}

	s.Require().NoError(p.buildValues())
}

func (s *BuildValuesTestSuite) TestValuesBuildError() {
	tmpDir := s.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	tmpValues := filepath.Join(tmpDir, "blablavalues.yaml")
	s.Require().NoError(os.WriteFile(tmpValues, []byte("a: b"), 0o600))

	mockedRelease := &mockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")

	errBuildValues := errors.New("values build error")
	mockedRelease.On("BuildValues").Return(errBuildValues)

	p.body = &planBody{
		Releases: releaseConfigs{mockedRelease},
	}

	s.Require().ErrorIs(p.buildValues(), errBuildValues)
	mockedRelease.AssertExpectations(s.T())
}

func (s *BuildValuesTestSuite) TestSuccess() {
	tmpDir := s.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	valuesName := "blablavalues.yaml"
	valuesContents := []byte("a: b")
	tmpValues := filepath.Join(tmpDir, valuesName)
	s.Require().NoError(os.WriteFile(tmpValues, valuesContents, 0o600))

	mockedRelease := &mockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Values").Return([]release.ValuesReference{
		{Src: tmpValues},
	})
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("BuildValues").Return(nil)

	p.body = &planBody{
		Releases: releaseConfigs{mockedRelease},
	}

	s.Require().NoError(p.buildValues())
	mockedRelease.AssertExpectations(s.T())
}

func TestBuildValuesTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildValuesTestSuite))
}
