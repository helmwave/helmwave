package repo_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/repo"
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
	schema := repo.Configs{}.JSONSchema()

	ts.Require().NotNil(schema)
	ts.Require().Equal("array", schema.Type)
}

func (ts *InterfaceTestSuite) TestUnmarshalYAMLEmpty() {
	var cfgs repo.Configs
	str := `[]`

	err := yaml.Unmarshal([]byte(str), &cfgs)

	ts.Require().NoError(err)
	ts.Require().Empty(cfgs)
}

func (ts *InterfaceTestSuite) TestUnmarshalYAML() {
	var cfgs repo.Configs
	str := `[{"name": "blabla"}]`

	err := yaml.Unmarshal([]byte(str), &cfgs)

	ts.Require().NoError(err)
	ts.Require().Len(cfgs, 1)
	ts.Require().Equal("blabla", cfgs[0].Name())
}
