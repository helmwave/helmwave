package registry_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/registry"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type InterfaceTestSuite struct {
	suite.Suite
}

func TestInterfaceTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(InterfaceTestSuite))
}

func (ts *InterfaceTestSuite) TestConfigsJSONSchema() {
	schema := registry.Configs{}.JSONSchema()

	ts.Require().NotNil(schema)
	ts.Require().Equal("array", schema.Type)
}

func (ts *InterfaceTestSuite) TestUnmarshalYAMLEmpty() {
	var cfgs registry.Configs
	str := `[]`

	err := yaml.Unmarshal([]byte(str), &cfgs)

	ts.Require().NoError(err)
	ts.Require().Empty(cfgs)
}

func (ts *InterfaceTestSuite) TestUnmarshalYAML() {
	var cfgs registry.Configs
	str := `[{"host": "blabla"}]`

	err := yaml.Unmarshal([]byte(str), &cfgs)

	ts.Require().NoError(err)
	ts.Require().Len(cfgs, 1)
	ts.Require().Equal(cfgs[0].Host(), "blabla")
}

func (ts *InterfaceTestSuite) TestUnmarshalYAMLInvalid() {
	var cfgs registry.Configs
	str := `{}`

	err := yaml.Unmarshal([]byte(str), &cfgs)

	ts.Require().ErrorIs(err, registry.YAMLDecodeError{})
}
