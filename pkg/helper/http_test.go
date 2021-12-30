package helper_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/stretchr/testify/suite"
)

type HTTPTestSuite struct {
	suite.Suite
}

func (s *HTTPTestSuite) TestDownloadBadURL() {
	tmpDir := s.T().TempDir()
	filePath := filepath.Join(tmpDir, "test")

	err := helper.Download(filePath, "\\asd://null")
	s.Require().Error(err)
}

func (s *HTTPTestSuite) TestDownloadBadResponse() {
	tmpDir := s.T().TempDir()
	filePath := filepath.Join(tmpDir, "test")
	url := fmt.Sprintf("https://helmwave.github.io/%s", s.T().Name())

	err := helper.Download(filePath, url)
	s.Require().Error(err)
}

func (s *HTTPTestSuite) TestDownload() {
	tmpDir := s.T().TempDir()
	filePath := filepath.Join(tmpDir, "test")
	url := "https://helmwave.github.io/"

	err := helper.Download(filePath, url)
	s.Require().NoError(err)
}

func (s *HTTPTestSuite) TestIsURL() {
	urls := []string{
		"https://blog.golang.org/slices-intro",
		"https://helmwave.github.io/",
	}

	for _, url := range urls {
		b := helper.IsURL(url)
		s.True(b)
	}
}

func TestHTTPTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HTTPTestSuite))
}
