package release_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type ValuesTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestValuesTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValuesTestSuite))
}

func (ts *ValuesTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *ValuesTestSuite) TestProhibitDst() {
	type config struct {
		Values []release.ValuesReference
	}

	src := `
values:
- src: 1
  dst: a
- src: 2
  dst: b
`
	c := &config{}

	err := yaml.Unmarshal([]byte(src), c)
	ts.Require().NoError(err)

	err = release.ProhibitDst(c.Values)
	ts.Require().Error(err)
}

func (ts *ValuesTestSuite) TestList() {
	type config struct {
		Values []release.ValuesReference
	}

	src := `
values:
- a
- b
`
	c := &config{}

	err := yaml.Unmarshal([]byte(src), c)
	ts.Require().NoError(err)

	ts.Require().Equal(&config{
		Values: []release.ValuesReference{
			{Src: "a"},
			{Src: "b"},
		},
	}, c)
}

func (ts *ValuesTestSuite) TestMap() {
	type config struct {
		Values []release.ValuesReference
	}

	src := `
values:
- src: 1
  render: false
- src: 2
  strict: true
`
	c := &config{}

	err := yaml.Unmarshal([]byte(src), c)
	ts.Require().NoError(err)

	ts.Require().Equal(&config{
		Values: []release.ValuesReference{
			{Src: "1", Strict: false},
			{Src: "2", Strict: true},
		},
	}, c)
}

func (ts *ValuesTestSuite) TestBuildNonExistingNonStrict() {
	r := release.NewConfig()
	r.ValuesF = []release.ValuesReference{
		{
			Src:    "nonexisting.values",
			Strict: false,
		},
	}

	_, err := r.BuildValues(ts.ctx, ".", template.TemplaterSprig, nil)

	ts.Require().NoError(err)
	ts.Require().Empty(r.Values())
}

func (ts *ValuesTestSuite) TestBuildNonExistingStrict() {
	r := release.NewConfig()
	r.ValuesF = []release.ValuesReference{
		{
			Src:    "nonexisting.values",
			Strict: true,
		},
	}

	_, err := r.BuildValues(ts.ctx, ".", template.TemplaterSprig, nil)

	ts.Require().Error(err)
}

func (ts *ValuesTestSuite) TestJSONSchema() {
	schema := (&release.ValuesReference{}).JSONSchema()

	ts.Require().NotNil(schema)

	ts.NotNil(schema.Properties.GetPair("src"))
	ts.NotNil(schema.Properties.GetPair("dst"))
	ts.NotNil(schema.Properties.GetPair("delimiter_left"))
	ts.NotNil(schema.Properties.GetPair("delimiter_right"))
	ts.NotNil(schema.Properties.GetPair("strict"))
	ts.NotNil(schema.Properties.GetPair("renderer"))
}

func (ts *ValuesTestSuite) TestGetValuesReleaseLevel() {
	tmpDir := ts.T().TempDir()

	values1 := filepath.Join(tmpDir, "values1.yaml")
	err := os.WriteFile(values1, []byte(`service:
  port: 8080`), 0o600)
	ts.Require().NoError(err)

	// values2: uses 1-arg getValues (current release)
	// values3: uses 2-arg getValues with current release name
	// values4: uses 2-arg getValues to forward to plan-level
	values2 := filepath.Join(tmpDir, "values2.yaml.tpl")
	err = os.WriteFile(values2, []byte(`{{ $v1 := getValues "`+values1+`" }}replica: {{ $v1.service.port }}`), 0o600)
	ts.Require().NoError(err)

	values3 := filepath.Join(tmpDir, "values3.yaml.tpl")
	err = os.WriteFile(values3, []byte(`{{ $v1 := getValues "app@default" "`+values1+`" }}count: {{ $v1.service.port }}`), 0o600)
	ts.Require().NoError(err)

	values4 := filepath.Join(tmpDir, "values4.yaml.tpl")
	err = os.WriteFile(values4, []byte(`{{ $redis := getValues "redis@default" "config.yaml" }}redis_host: {{ $redis.host }}`), 0o600)
	ts.Require().NoError(err)

	r := release.NewConfig()
	r.NameF = "app"
	r.NamespaceF = "default"
	r.ValuesF = []release.ValuesReference{
		{Src: values1, Renderer: "copy"},
		{Src: values2, Renderer: "sprig"},
		{Src: values3, Renderer: "sprig"},
		{Src: values4, Renderer: "sprig"},
	}

	templateFuncs := map[string]any{
		"getValues": func(rel string, filename string) (any, error) {
			if rel == "redis@default" && filename == "config.yaml" {
				return map[string]any{"host": "redis.example.com"}, nil
			}

			return nil, fmt.Errorf("release %q file %q not found in plan", rel, filename)
		},
	}

	renderedValues, err := r.BuildValues(ts.ctx, tmpDir, template.TemplaterSprig, templateFuncs)
	ts.Require().NoError(err)

	ts.Require().Contains(renderedValues[values2], "replica: 8080")
	ts.Require().Contains(renderedValues[values3], "count: 8080")
	ts.Require().Contains(renderedValues[values4], "redis_host: redis.example.com")
}
