package helper_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/stretchr/testify/suite"
)

type HTTPTestSuite struct {
	suite.Suite
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
