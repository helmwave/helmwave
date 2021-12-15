//go:build ignore || unit

package helper

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type HelmTestSuite struct {
	suite.Suite
}

func (s *HelmTestSuite) TestHelmNS() {
	h1, err := NewHelm("my")

	s.Require().NoError(err)
	s.Require().Equal("my", h1.Namespace())
}

func TestHelmTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HelmTestSuite))
}
