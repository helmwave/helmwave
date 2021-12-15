//go:build ignore || unit

package release

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestConfigUniq() {
	r := &Config{
		Name:      "redis",
		Namespace: "test",
	}

	s.Require().Equal(r.Uniq(), r.uniqName)
	s.Require().True(r.Uniq().Validate())
}

func TestConfigTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ConfigTestSuite))
}
