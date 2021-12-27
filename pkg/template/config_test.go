package template_test

import (
	"net/url"
	"testing"

	"github.com/helmwave/helmwave/pkg/template"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestUnmarshal() {
	u, err := url.Parse("http://blabla")
	s.Require().NoError(err)

	source := &template.Config{}
	source.Gomplate.Datasources = map[string]template.Source{
		"test": {
			u,
		},
	}
	dest := &template.Config{
		Gomplate: template.GomplateConfig{},
	}

	data, err := yaml.Marshal(source)
	s.Require().NoError(err)

	s.Require().NoError(yaml.Unmarshal(data, dest))

	s.Require().Equal(source, dest)
}

func TestConfigTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ConfigTestSuite))
}
