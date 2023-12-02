package hooks_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/hooks"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type YAMLTestSuite struct {
	suite.Suite
}

func TestYAMLTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(YAMLTestSuite))
}

func (s *YAMLTestSuite) TestEmptyStructure() {
	lifecycle := hooks.Lifecycle{}
	str := `{}`

	err := yaml.Unmarshal([]byte(str), &lifecycle)

	s.Require().NoError(err)

	s.Require().Empty(lifecycle.PreBuild)
	s.Require().Empty(lifecycle.PostBuild)
	s.Require().Empty(lifecycle.PreUp)
	s.Require().Empty(lifecycle.PostUp)
	s.Require().Empty(lifecycle.PreRollback)
	s.Require().Empty(lifecycle.PostRollback)
	s.Require().Empty(lifecycle.PreDown)
	s.Require().Empty(lifecycle.PostDown)
}

func (s *YAMLTestSuite) TestHooksUnmarshal() {
	lifecycle := hooks.Lifecycle{}
	str := `
pre_build:
  - test123
`

	err := yaml.Unmarshal([]byte(str), &lifecycle)

	s.Require().NoError(err)

	s.Require().Len(lifecycle.PreBuild, 1)
	s.Require().Empty(lifecycle.PostBuild)
	s.Require().Empty(lifecycle.PreUp)
	s.Require().Empty(lifecycle.PostUp)
	s.Require().Empty(lifecycle.PreRollback)
	s.Require().Empty(lifecycle.PostRollback)
	s.Require().Empty(lifecycle.PreDown)
	s.Require().Empty(lifecycle.PostDown)

	hook := lifecycle.PreBuild[0]
	s.Require().NotNil(hook)
	s.Require().IsType(hooks.NewHook(), hook)
}

func (s *YAMLTestSuite) TestUnmarshalShortForm() {
	hook := hooks.NewHook()
	str := `test 123 456`

	err := yaml.Unmarshal([]byte(str), hook)

	s.Require().NoError(err)
	s.Require().Equal("test", hook.Cmd)
	s.Require().Equal([]string{"123", "456"}, hook.Args)
	s.Require().False(hook.AllowFailure)
	s.Require().True(hook.Show)
}

func (s *YAMLTestSuite) TestUnmarshalLongForm() {
	hook := hooks.NewHook()
	str := `
cmd: test 123
args:
  - 456
show: false
allow_failure: true
`

	err := yaml.Unmarshal([]byte(str), hook)

	s.Require().NoError(err)
	s.Require().Equal("test 123", hook.Cmd)
	s.Require().Equal([]string{"456"}, hook.Args)
	s.Require().True(hook.AllowFailure)
	s.Require().False(hook.Show)
}
