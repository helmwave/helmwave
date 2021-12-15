//go:build ignore || unit

package template

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestUnmarshal() {
	u, err := url.Parse("http://blabla")
	s.Require().NoError(err)

	source := &Config{
		Gomplate: GomplateConfig{
			Datasources: map[string]Source{
				"test": Source{
					u,
				},
			},
		},
	}
	dest := &Config{
		Gomplate: GomplateConfig{},
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
