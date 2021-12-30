package helper_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/stretchr/testify/suite"
)

type HelmTestSuite struct {
	suite.Suite
}

func (s *HelmTestSuite) TestNewCfg() {
	ns := s.T().Name()
	cfg, err := helper.NewCfg(ns)

	s.Require().NoError(err)
	s.Require().NotNil(cfg)
}

func (s *HelmTestSuite) TestNewHelmNS() {
	ns := s.T().Name()
	h1, err := helper.NewHelm(ns)

	s.Require().NoError(err)
	s.Require().NotNil(h1)
	s.Require().Equal(ns, h1.Namespace())
}

func TestHelmTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HelmTestSuite))
}
