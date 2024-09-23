package release_test

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/fileref"

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
		Values []fileref.Config
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

	err = fileref.ProhibitDst(c.Values)
	ts.Require().Error(err)
}

func (ts *ValuesTestSuite) TestList() {
	type config struct {
		Values []fileref.Config
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
		Values: []fileref.Config{
			{Src: "a"},
			{Src: "b"},
		},
	}, c)
}

func (ts *ValuesTestSuite) TestMap() {
	type config struct {
		Values []fileref.Config
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
		Values: []fileref.Config{
			{Src: "1", Strict: false},
			{Src: "2", Strict: true},
		},
	}, c)
}

func (ts *ValuesTestSuite) TestBuildNonExistingNonStrict() {
	r := release.NewConfig()
	r.ValuesF = []fileref.Config{
		{
			Src:    "non-existing.values",
			Strict: false,
		},
	}

	err := r.BuildValues(ts.ctx, ".", template.TemplaterSprig)

	ts.Require().NoError(err)
	ts.Require().Empty(r.Values())
}

func (ts *ValuesTestSuite) TestBuildNonExistingStrict() {
	r := release.NewConfig()
	r.ValuesF = []fileref.Config{
		{
			Src:    "non-existing.values",
			Strict: true,
		},
	}

	err := r.BuildValues(ts.ctx, ".", template.TemplaterSprig)

	ts.Require().Error(err)
}

func (ts *ValuesTestSuite) TestJSONSchema() {
	schema := (&fileref.Config{}).JSONSchema()

	ts.Require().NotNil(schema)

	ts.NotNil(schema.Properties.GetPair("src"))
	ts.NotNil(schema.Properties.GetPair("dst"))
	ts.NotNil(schema.Properties.GetPair("delimiter_left"))
	ts.NotNil(schema.Properties.GetPair("delimiter_right"))
	ts.NotNil(schema.Properties.GetPair("strict"))
	ts.NotNil(schema.Properties.GetPair("renderer"))
}
