//go:build ignore || integration

package action

import (
	"os"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type BuildTestSuite struct {
	suite.Suite
}

func (ts *BuildTestSuite) TearDownTest() {
	_ = os.RemoveAll(tests.Root + plan.Dir)
}

func (ts *BuildTestSuite) BuildManifest() {
	y := &Yml{
		tests.Root + "01_helmwave.yml.tpl",
		tests.Root + "02_helmwave.yml",
	}

	s := &Build{
		plandir:  tests.Root + plan.Dir,
		yml:      y,
		tags:     cli.StringSlice{},
		matchAll: true,
	}

	ts.Require().NoError(s.Run())
	ts.Require().DirExists(tests.Root + plan.Dir + plan.Manifest)
}

// func (ts *BuildTestSuite) BuildRepositories404() {
//	s := &Build{
//		plandir:  tests.Root + plan.Dir,
//		ymlFile:      tests.Root + "04_helmwave.yml",
//		tags:     cli.StringSlice{},
//		matchAll: true,
//	}
//
//	err := s.Run()
//	if !errors.Is(err, repo.ErrNotFound) && err != nil {
//		t.Error("'bitnami' must be not found")
//	}
// }

func (ts *BuildTestSuite) BuildRepositories() {
	y := &Yml{
		tests.Root + "01_helmwave.yml.tpl",
		tests.Root + "02_helmwave.yml",
	}

	s := &Build{
		plandir:  tests.Root + plan.Dir,
		yml:      y,
		tags:     cli.StringSlice{},
		matchAll: true,
	}

	ts.Require().NoError(s.Run())

	const rep = "bitnami"
	b, _ := plan.NewBody(tests.Root + plan.Dir + plan.File)

	if _, found := repo.IndexOfName(b.Repositories, rep); !found {
		ts.Failf("%q not found", rep)
	}
}

func (ts *BuildTestSuite) BuildReleasesMatchGroup() {
	y := &Yml{
		tests.Root + "01_helmwave.yml.tpl",
		tests.Root + "03_helmwave.yml",
	}

	cases := []struct {
		tags  *cli.StringSlice
		names []string
	}{
		{
			tags:  cli.NewStringSlice("b"),
			names: []string{"redis-b", "memcached-b"},
		},
		{
			tags:  cli.NewStringSlice("b", "redis"),
			names: []string{"redis-b"},
		},
	}

	for _, c := range cases {
		s := &Build{
			plandir:  tests.Root + plan.Dir,
			yml:      y,
			tags:     *c.tags,
			matchAll: true,
		}

		ts.Require().NoError(s.Run())

		b, _ := plan.NewBody(tests.Root + plan.Dir + plan.File)

		names := make([]string, 0, len(b.Releases))
		for _, r := range b.Releases {
			names = append(names, r.Name)
		}

		ts.Require().ElementsMatch(c.names, names)
	}

}

func (ts *BuildTestSuite) BuildAutoYml() {
	y := &Yml{
		tests.Root + "01_helmwave.yml.tpl",
		tests.Root + "01_auto_yaml_helmwave.yml",
	}

	s := &Build{
		plandir:  tests.Root + plan.Dir,
		tags:     cli.StringSlice{},
		matchAll: true,
		autoYml:  true,
		yml:      y,
	}

	value := "test01"
	_ = os.Setenv("PROJECT_NAME", value)
	_ = os.Setenv("NAMESPACE", value)

	ts.Require().NoError(s.Run())
	ts.Require().DirExists(tests.Root + plan.Dir + plan.Manifest)
}

func (ts *BuildTestSuite) BuildGomplate() {
	y := &Yml{
		tests.Root + "08_helmwave.yml",
		tests.Root + "08_values.yml",
	}

	s := &Build{
		plandir:  tests.Root + plan.Dir,
		tags:     cli.StringSlice{},
		matchAll: true,
		autoYml:  true,
		yml:      y,
	}

	value := "test08"
	_ = os.Setenv("PROJECT_NAME", value)
	_ = os.Setenv("NAMESPACE", value)

	ts.Require().NoError(s.Run())
	ts.Require().DirExists(tests.Root + plan.Dir + plan.Manifest)
}

func TestBuildTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(BuildTestSuite))
}
