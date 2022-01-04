package release_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestConfigUniq() {
	r := release.NewConfig()
	r.NameF = "redis"
	r.NamespaceF = "test"

	s.Require().NoError(r.Uniq().Validate())
}

func TestConfigTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ConfigTestSuite))
}
