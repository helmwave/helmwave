package helper_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type HTTPTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestHTTPTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HTTPTestSuite))
}

func (ts *HTTPTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *HTTPTestSuite) TestDownloadBadURL() {
	tmpDir := ts.T().TempDir()
	filePath := filepath.Join(tmpDir, "test")

	err := helper.Download(ts.ctx, filePath, "\\asd://null")
	ts.Require().Error(err)
}

func (ts *HTTPTestSuite) TestDownloadBadResponse() {
	tmpDir := ts.T().TempDir()
	filePath := filepath.Join(tmpDir, "test")
	url := fmt.Sprintf("https://helmwave.github.io/%s", ts.T().Name())

	err := helper.Download(ts.ctx, filePath, url)
	ts.Require().Error(err)
}

func (ts *HTTPTestSuite) TestDownload() {
	tmpDir := ts.T().TempDir()
	filePath := filepath.Join(tmpDir, "test")
	url := "https://helmwave.github.io/"

	err := helper.Download(ts.ctx, filePath, url)
	ts.Require().NoError(err)
}

func (ts *HTTPTestSuite) TestIsURL() {
	urls := []string{
		"https://blog.golang.org/slices-intro",
		"https://helmwave.github.io/",
	}

	for _, url := range urls {
		b := helper.IsURL(url)
		ts.True(b)
	}
}
