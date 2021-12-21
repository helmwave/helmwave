package helper_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/stretchr/testify/suite"
)

type HelmTestSuite struct {
	suite.Suite
}

func (s *HelmTestSuite) TestHelmNS() {
	h1, err := helper.NewHelm("my")

	s.Require().NoError(err)
	s.Require().Equal("my", h1.Namespace())
}

func TestHelmTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HelmTestSuite))
}
