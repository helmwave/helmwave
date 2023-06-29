package hooks_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/hooks"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) TestHooksJSONSchema() {
	schema := hooks.Hooks{}.JSONSchema()

	s.Require().NotNil(schema)
	s.Require().Equal("array", schema.Type)
}

func (s *ConfigTestSuite) TestHookJSONSchema() {
	schema := hooks.NewHook().JSONSchema()

	s.Require().NotNil(schema)

	keys := schema.Properties.Keys()
	s.Require().Contains(keys, "cmd")
	s.Require().Contains(keys, "args")
	s.Require().Contains(keys, "show")
	s.Require().Contains(keys, "allow_failure")
}
