//go:build ignore || unit

package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ExtraTestSuite struct {
	suite.Suite
}

func (s *ExtraTestSuite) TestUnmarshal() {
	s.T().Skip("not implemented")
}

func (s *ExtraTestSuite) TestToYaml() {
	s.T().Skip("not implemented")
}

func (s *ExtraTestSuite) TestFromYaml() {
	s.T().Skip("not implemented")
}

func (s *ExtraTestSuite) TestExec() {
	s.T().Skip("not implemented")
}

func (s *ExtraTestSuite) TestSetValueAtPath() {
	s.T().Skip("not implemented")
}

func (s *ExtraTestSuite) TestRequiredEnv() {
	s.T().Skip("not implemented")
}

func (s *ExtraTestSuite) TestRequired() {
	tests := []struct {
		data  interface{}
		fails bool
	}{
		{
			data:  nil,
			fails: true,
		},
		{
			data:  4,
			fails: false,
		},
		{
			data:  "",
			fails: true,
		},
		{
			data:  "123",
			fails: false,
		},
	}

	for _, t := range tests {
		res, err := Required("blabla", t.data)
		if t.fails {
			s.Require().Error(err)
			s.Require().Nil(res)
		} else {
			s.Require().NoError(err)
			s.Require().Equal(t.data, res)
		}
	}
}

func (s *ExtraTestSuite) TestReadFile() {
	tmpDir := s.T().TempDir()
	tmpFile := filepath.Join(tmpDir, "blablafile")

	res, err := ReadFile(tmpFile)

	s.Require().Equal("", res)
	s.Require().ErrorIs(err, os.ErrNotExist)

	data := s.T().Name()

	s.Require().NoError(os.WriteFile(tmpFile, []byte(data), 0666))
	s.Require().FileExists(tmpFile)

	res, err = ReadFile(tmpFile)

	s.Require().NoError(err)
	s.Require().Equal(data, res)
}

func (s *ExtraTestSuite) TestGet() {
	s.T().Skip("not implemented")
}

func (s *ExtraTestSuite) TestHasKey() {
	s.T().Skip("not implemented")
}

func TestExtraTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExtraTestSuite))
}
