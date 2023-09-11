package monitor_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/registry"
	"github.com/stretchr/testify/suite"
)

type InterfaceTestSuite struct {
	suite.Suite
}

func TestInterfaceTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(InterfaceTestSuite))
}

func (s *InterfaceTestSuite) TestConfigsJSONSchema() {
	schema := registry.Configs{}.JSONSchema()

	s.Require().NotNil(schema)
	s.Require().Equal("array", schema.Type)
}
