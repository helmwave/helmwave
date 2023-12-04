package hooks_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/hooks"
	"github.com/invopop/jsonschema"
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
	reflector := &jsonschema.Reflector{DoNotReference: true}
	schema := reflector.Reflect(hooks.NewHook())

	s.Require().NotNil(schema)

	s.NotNil(schema.Properties.GetPair("cmd"))
	s.NotNil(schema.Properties.GetPair("args"))
	s.NotNil(schema.Properties.GetPair("show"))
	s.NotNil(schema.Properties.GetPair("allow_failure"))
}
