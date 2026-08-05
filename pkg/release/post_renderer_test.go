package release_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type PostRendererTestSuite struct {
	suite.Suite
}

func TestPostRendererTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(PostRendererTestSuite))
}

func (ts *PostRendererTestSuite) TestUnmarshalLegacyArray() {
	var cfg release.PostRendererReference
	str := `["./script.sh", "--arg1", "--arg2"]`

	err := yaml.Unmarshal([]byte(str), &cfg)

	ts.Require().NoError(err)
	ts.Require().Equal(release.ExecConfig{"./script.sh", "--arg1", "--arg2"}, cfg.Exec)
	ts.Require().Nil(cfg.Kustomize)
}

func (ts *PostRendererTestSuite) TestUnmarshalExecString() {
	var cfg release.PostRendererReference
	str := `exec: "./script.sh --arg1 --arg2"`

	err := yaml.Unmarshal([]byte(str), &cfg)

	ts.Require().NoError(err)
	ts.Require().Equal(release.ExecConfig{"./script.sh", "--arg1", "--arg2"}, cfg.Exec)
	ts.Require().Nil(cfg.Kustomize)
}

func (ts *PostRendererTestSuite) TestUnmarshalExecArray() {
	var cfg release.PostRendererReference
	str := `exec: ["./script.sh", "--arg1", "--arg2"]`

	err := yaml.Unmarshal([]byte(str), &cfg)

	ts.Require().NoError(err)
	ts.Require().Equal(release.ExecConfig{"./script.sh", "--arg1", "--arg2"}, cfg.Exec)
	ts.Require().Nil(cfg.Kustomize)
}

func (ts *PostRendererTestSuite) TestUnmarshalKustomizeString() {
	var cfg release.PostRendererReference
	str := `kustomize: ./my-kustomize`

	err := yaml.Unmarshal([]byte(str), &cfg)

	ts.Require().NoError(err)
	ts.Require().Empty(cfg.Exec)
	ts.Require().NotNil(cfg.Kustomize)
	ts.Require().Equal("./my-kustomize", cfg.Kustomize.Src)
}

func (ts *PostRendererTestSuite) TestUnmarshalKustomizeObject() {
	var cfg release.PostRendererReference
	str := `
kustomize:
  src: ./my-kustomize
  renderer: gomplate
  delimiter_left: "[["
  delimiter_right: "]]"
`

	err := yaml.Unmarshal([]byte(str), &cfg)

	ts.Require().NoError(err)
	ts.Require().Empty(cfg.Exec)
	ts.Require().NotNil(cfg.Kustomize)
	ts.Require().Equal("./my-kustomize", cfg.Kustomize.Src)
	ts.Require().Equal("gomplate", cfg.Kustomize.Renderer)
	ts.Require().Equal("[[", cfg.Kustomize.DelimiterLeft)
	ts.Require().Equal("]]", cfg.Kustomize.DelimiterRight)
}

func (ts *PostRendererTestSuite) TestUnmarshalMutualExclusivity() {
	var cfg release.PostRendererReference
	str := `
exec: "./script.sh"
kustomize: ./my-kustomize
`

	err := yaml.Unmarshal([]byte(str), &cfg)

	ts.Require().Error(err)
	ts.Require().Contains(err.Error(), "mutually exclusive")
}

func (ts *PostRendererTestSuite) TestIsEmpty() {
	ts.True((&release.PostRendererReference{}).IsEmpty())
	ts.True((*release.PostRendererReference)(nil).IsEmpty())

	ts.False((&release.PostRendererReference{Exec: []string{"cmd"}}).IsEmpty())
	ts.False((&release.PostRendererReference{Kustomize: &release.KustomizeConfig{Src: "path"}}).IsEmpty())
}

func (ts *PostRendererTestSuite) TestHelmPostRendererExec() {
	cfg := &release.PostRendererReference{Exec: []string{"/bin/sh", "-c", "cat"}}

	pr, err := cfg.HelmPostRenderer()

	ts.Require().NoError(err)
	ts.Require().NotNil(pr)
}

func (ts *PostRendererTestSuite) TestHelmPostRendererKustomizeNotBuilt() {
	cfg := &release.PostRendererReference{Kustomize: &release.KustomizeConfig{Src: "./kustomize"}}

	_, err := cfg.HelmPostRenderer()

	ts.Require().Error(err)
	ts.Require().Contains(err.Error(), "not built yet")
}

func (ts *PostRendererTestSuite) TestJSONSchema() {
	schema := release.PostRendererReference{}.JSONSchema()

	ts.Require().NotNil(schema)
	ts.Require().NotNil(schema.OneOf)
	ts.Require().Len(schema.OneOf, 2)
}
