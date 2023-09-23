package helper_test

import (
	"io/fs"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/stretchr/testify/suite"
)

type FileTestSuite struct {
	suite.Suite
}

func (s *FileTestSuite) TestContainsGood() {
	b := helper.Contains("c", []string{
		"a",
		"b",
		"c",
		"d",
		"c",
	})

	s.Require().True(b)
}

func (s *FileTestSuite) TestContainsBad() {
	b := helper.Contains("12", []string{
		"a",
		"b",
		"c",
		"d",
		"c",
	})

	s.Require().False(b)
}

func (s *FileTestSuite) TestCreateFile() {
	tmpDir := s.T().TempDir()
	tmpFS, _ := filefs.New(&url.URL{Scheme: "file", Path: tmpDir})
	filePath := filepath.Join("testdir", "test")

	f, err := helper.CreateFile(tmpFS.(fsimpl.WriteableFS), filePath)
	s.Require().NoError(err)
	s.Require().NotNil(f)
	s.Require().NoError(f.Close())

	s.Require().FileExists(filepath.Join(tmpDir, filePath))
}

func (s *FileTestSuite) TestCreateFileMkdir() {
	tmpDir := s.T().TempDir()
	tmpFS, _ := filefs.New(&url.URL{Scheme: "file", Path: tmpDir})

	f, err := helper.CreateFile(tmpFS.(fsimpl.WriteableFS), ".")
	s.Require().Error(err)
	s.Require().Nil(f)

	filePath := "test"

	f, err = helper.CreateFile(tmpFS.(fsimpl.WriteableFS), filePath)
	s.Require().NoError(err)
	s.Require().NotNil(f)
	s.Require().NoError(f.Close())

	filePath = filepath.Join(filePath, "test")

	f, err = helper.CreateFile(tmpFS.(fsimpl.WriteableFS), filePath)
	s.Require().Error(err)
	s.Require().Nil(f)
}

func (s *FileTestSuite) TestIsExists() {
	tmpDir := s.T().TempDir()
	tmpFS, _ := filefs.New(&url.URL{Scheme: "file", Path: tmpDir})

	filePath := "test"

	s.Require().False(helper.IsExists(tmpFS.(fs.StatFS), filePath))

	f, err := helper.CreateFile(tmpFS.(fsimpl.WriteableFS), filePath)
	s.Require().NoError(err)
	s.Require().NotNil(f)
	s.Require().NoError(f.Close())

	s.Require().True(helper.IsExists(tmpFS.(fs.StatFS), filePath))
}

func TestFileTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(FileTestSuite))
}
