package helper_test

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/stretchr/testify/suite"
)

type HTTPTestSuite struct {
	suite.Suite
}

func (s *HTTPTestSuite) TestDownloadBadURL() {
	tmpDir := s.T().TempDir()
	filePath := filepath.Join(tmpDir, "test")

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := helper.Download(baseFS.(fsimpl.WriteableFS), filePath, "\\asd://null")
	s.Require().Error(err)
}

func (s *HTTPTestSuite) TestDownloadBadResponse() {
	tmpDir := s.T().TempDir()
	filePath := filepath.Join(tmpDir, "test")
	u := fmt.Sprintf("https://helmwave.github.io/%s", s.T().Name())

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := helper.Download(baseFS.(fsimpl.WriteableFS), filePath, u)
	s.Require().Error(err)
}

func (s *HTTPTestSuite) TestDownload() {
	tmpDir := s.T().TempDir()
	filePath := filepath.Join(tmpDir, "test")
	u := "https://helmwave.github.io/"

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := helper.Download(baseFS.(fsimpl.WriteableFS), filePath, u)
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
