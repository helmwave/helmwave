//go:build ignore || unit

package helper

import (
	"testing"

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
		b := IsURL(url)
		s.True(b)
	}
}

func TestHTTPTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HTTPTestSuite))
}
