//go:build ignore || unit

package action

import (
	"path/filepath"
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

func (ts *BuildTestSuite) TestManifest() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
		filepath.Join(tests.Root, "02_helmwave.yml"),
	}

	s := &Build{
		plandir:  tmpDir,
		yml:      y,
		tags:     cli.StringSlice{},
		matchAll: true,
	}

	ts.Require().NoError(s.Run())
	ts.Require().DirExists(filepath.Join(s.plandir, plan.Manifest))
}

// func (ts *BuildTestSuite) TestRepositories404() {
//	s := &Build{
//		plandir:  tmpDir,
//		ymlFile:      filepath.Join(tests.Root, "04_helmwave.yml"),
//		tags:     cli.StringSlice{},
//		matchAll: true,
//	}
//
//	err := s.Run()
//	if !errors.Is(err, repo.ErrNotFound) && err != nil {
//		t.Error("'bitnami' must be not found")
//	}
// }

func (ts *BuildTestSuite) TestRepositories() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
		filepath.Join(tests.Root, "02_helmwave.yml"),
	}

	s := &Build{
		plandir:  tmpDir,
		yml:      y,
		tags:     cli.StringSlice{},
		matchAll: true,
	}

	ts.Require().NoError(s.Run())

	const rep = "bitnami"
	b, _ := plan.NewBody(filepath.Join(s.plandir, plan.File))

	if _, found := repo.IndexOfName(b.Repositories, rep); !found {
		ts.Failf("%q not found", rep)
	}
}

func (ts *BuildTestSuite) TestReleasesMatchGroup() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
		filepath.Join(tests.Root, "03_helmwave.yml"),
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
			plandir:  tmpDir,
			yml:      y,
			tags:     *c.tags,
			matchAll: true,
		}

		ts.Require().NoError(s.Run())

		b, _ := plan.NewBody(filepath.Join(s.plandir, plan.File))

		names := make([]string, 0, len(b.Releases))
		for _, r := range b.Releases {
			names = append(names, r.Name)
		}

		ts.Require().ElementsMatch(c.names, names)
	}

}

func (ts *BuildTestSuite) TestAutoYml() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
		filepath.Join(tmpDir, "01_auto_yaml_helmwave.yml"),
	}

	s := &Build{
		plandir:  tmpDir,
		tags:     cli.StringSlice{},
		matchAll: true,
		autoYml:  true,
		yml:      y,
	}

	value := "test01"
	ts.T().Setenv("PROJECT_NAME", value)
	ts.T().Setenv("NAMESPACE", value)

	ts.Require().NoError(s.Run())
	ts.Require().DirExists(filepath.Join(s.plandir, plan.Manifest))
}

func (ts *BuildTestSuite) TestGomplate() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		filepath.Join(tests.Root, "08_helmwave.yml"),
		filepath.Join(tmpDir, "08_helmwave.yml"),
	}

	s := &Build{
		plandir:  tmpDir,
		tags:     cli.StringSlice{},
		matchAll: true,
		autoYml:  true,
		yml:      y,
	}

	value := "test08"
	ts.T().Setenv("PROJECT_NAME", value)
	ts.T().Setenv("NAMESPACE", value)

	ts.Require().NoError(s.Run())
	ts.Require().DirExists(filepath.Join(s.plandir, plan.Manifest))
}

func (ts *BuildTestSuite) TestDiffLocal() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		filepath.Join(tests.Root, "08_helmwave.yml"),
		filepath.Join(tmpDir, "08_values.yml"),
	}

	s := &Build{
		plandir:  tmpDir,
		tags:     cli.StringSlice{},
		matchAll: true,
		autoYml:  true,
		yml:      y,
		diff:     &Diff{},
		diffMode: diffModeLocal,
	}

	ts.Require().NoError(s.Run(), "build should not fail without diffing")
	ts.Require().NoError(s.Run(), "build should not fail with diffing with previous plan")
}

func TestBuildTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildTestSuite))
}
