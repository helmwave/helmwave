package release_test

import (
	"io/fs"
	"net/url"
	"os"
	"testing"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type ValuesTestSuite struct {
	suite.Suite
}

func TestValuesTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValuesTestSuite))
}

func (s *ValuesTestSuite) TestList() {
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
	s.Require().NoError(err)

	s.Require().Equal(&config{
		Values: []release.ValuesReference{
			{Src: "a"},
			{Src: "b"},
		},
	}, c)
}

func (s *ValuesTestSuite) TestMap() {
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
	s.Require().NoError(err)

	s.Require().Equal(&config{
		Values: []release.ValuesReference{
			{Src: "1", Strict: false},
			{Src: "2", Strict: true},
		},
	}, c)
}

func (s *ValuesTestSuite) TestBuildNonExistingNonStrict() {
	r := release.NewConfig()
	r.ValuesF = []release.ValuesReference{
		{
			Src:    "nonexisting.values",
			Strict: false,
		},
	}

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := r.ExportValues(baseFS.(fs.StatFS), baseFS.(fsimpl.WriteableFS), template.TemplaterSprig)

	s.Require().NoError(err)
	s.Require().Len(r.Values(), 0)
}

func (s *ValuesTestSuite) TestBuildNonExistingStrict() {
	r := release.NewConfig()
	r.ValuesF = []release.ValuesReference{
		{
			Src:    "nonexisting.values",
			Strict: true,
		},
	}

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := r.ExportValues(baseFS.(fs.StatFS), baseFS.(fsimpl.WriteableFS), template.TemplaterSprig)

	s.Require().Error(err)
}

func (s *ValuesTestSuite) TestJSONSchema() {
	schema := (&release.ValuesReference{}).JSONSchema()

	s.Require().NotNil(schema)

	keys := schema.Properties.Keys()
	s.Require().Contains(keys, "src")
	s.Require().Contains(keys, "dst")
	s.Require().Contains(keys, "delimiter_left")
	s.Require().Contains(keys, "delimiter_right")
	s.Require().Contains(keys, "strict")
	s.Require().Contains(keys, "renderer")
}
