// +build ignore unit

package helper

import (
	helm "helm.sh/helm/v3/pkg/cli"
)

func (s *HelperTestSuite) TestActionCfg() {
	settings := helm.New()

	cfg, err := ActionCfg("", settings)
	s.Require().NoError(err)
	s.Require().NotNil(cfg)
	s.Require().NotNil(cfg.Releases)
}
