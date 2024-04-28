package helper_test

import (
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/stretchr/testify/suite"
)

type FileTestSuite struct {
	suite.Suite
}

func (s *FileTestSuite) TestCreateFile() {
	tmpDir := s.T().TempDir()
	filePath := filepath.Join(tmpDir, "testdir", "test")

	f, err := helper.CreateFile(filePath)
	s.Require().NoError(err)
	s.Require().NotNil(f)
	s.Require().NoError(f.Close())

	s.Require().FileExists(filePath)
}

func (s *FileTestSuite) TestCreateFileMkdir() {
	tmpDir := s.T().TempDir()

	f, err := helper.CreateFile(tmpDir)
	s.Require().Error(err)
	s.Require().Nil(f)

	filePath := filepath.Join(tmpDir, "test")

	f, err = helper.CreateFile(filePath)
	s.Require().NoError(err)
	s.Require().NotNil(f)
	s.Require().NoError(f.Close())

	filePath = filepath.Join(filePath, "test")

	f, err = helper.CreateFile(filePath)
	s.Require().Error(err)
	s.Require().Nil(f)
}

func (s *FileTestSuite) TestIsExists() {
	tmpDir := s.T().TempDir()
	filePath := filepath.Join(tmpDir, "test")

	s.Require().False(helper.IsExists(filePath))

	f, err := helper.CreateFile(filePath)
	s.Require().NoError(err)
	s.Require().NotNil(f)
	s.Require().NoError(f.Close())

	s.Require().True(helper.IsExists(filePath))
}

func TestFileTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(FileTestSuite))
}
