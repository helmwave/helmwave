package repo_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/invopop/jsonschema"
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

func (s *InterfaceTestSuite) TestConfigsJSONSchema() {
	schema := jsonschema.Reflect(repo.Configs{})

	s.Require().NotNil(schema)
	s.Require().Equal("array", schema.Type)
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
	ts.Require().Equal(cfgs[0].Name(), "blabla")
}
