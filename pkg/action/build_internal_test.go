package action

import (
	"path/filepath"
	"strings"
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

func (ts *BuildTestSuite) TestImplementsAction() {
	ts.Require().Implements((*Action)(nil), &Build{})
}

func (ts *BuildTestSuite) TestManifest() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
		file:      filepath.Join(tests.Root, "02_helmwave.yml"),
		templater: "sprig",
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
		tpl:       filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
		file:      filepath.Join(tests.Root, "02_helmwave.yml"),
		templater: "sprig",
	}

	s := &Build{
		plandir:  tmpDir,
		yml:      y,
		tags:     cli.StringSlice{},
		matchAll: true,
	}

	ts.Require().NoError(s.Run())

	const rep = "bitnami"
	b, err := plan.NewBody(filepath.Join(s.plandir, plan.File))

	if _, found := repo.IndexOfName(b.Repositories, rep); !found {
		ts.Failf("%q not found", rep)
	}
}

func (ts *BuildTestSuite) TestReleasesMatchGroup() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
		file:      filepath.Join(tests.Root, "03_helmwave.yml"),
		templater: "sprig",
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

	for i := range cases {
		s := &Build{
			plandir:  tmpDir,
			yml:      y,
			tags:     *cases[i].tags,
			matchAll: true,
		}

		ts.Require().NoError(s.Run())

		b, _ := plan.NewBody(filepath.Join(s.plandir, plan.File))

		names := make([]string, 0, len(b.Releases))
		for _, r := range b.Releases {
			names = append(names, r.Name())
		}

		ts.Require().ElementsMatch(cases[i].names, names)
	}
}

func (ts *BuildTestSuite) TestDiffLocal() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "07_helmwave.yml"),
		file:      filepath.Join(tests.Root, "07_helmwave.yml"),
		templater: "sprig",
	}

	s := &Build{
		plandir:  tmpDir,
		tags:     cli.StringSlice{},
		matchAll: true,
		autoYml:  true,
		yml:      y,
		diff:     &Diff{},
		diffMode: DiffModeLocal,
	}

	ts.Require().NoError(s.Run(), "build should not fail without diffing")
	ts.Require().NoError(s.Run(), "build should not fail with diffing with previous plan")
}

func TestBuildTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildTestSuite))
}

type NonParallelBuildTestSuite struct {
	suite.Suite
}

func (ts *NonParallelBuildTestSuite) TestAutoYml() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
		file:      filepath.Join(tmpDir, "01_auto_yaml_helmwave.yml"),
		templater: "sprig",
	}

	s := &Build{
		plandir:  tmpDir,
		tags:     cli.StringSlice{},
		matchAll: true,
		autoYml:  true,
		yml:      y,
	}

	value := strings.ToLower(strings.ReplaceAll(ts.T().Name(), "/", ""))
	ts.T().Setenv("NAMESPACE", value)

	ts.Require().NoError(s.Run())
	ts.Require().DirExists(filepath.Join(s.plandir, plan.Manifest))
}

func (ts *NonParallelBuildTestSuite) TestGomplate() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "08_helmwave.yml"),
		file:      filepath.Join(tmpDir, "08_helmwave.yml"),
		templater: "gomplate",
	}

	s := &Build{
		plandir:  tmpDir,
		tags:     cli.StringSlice{},
		matchAll: true,
		autoYml:  true,
		yml:      y,
	}

	ts.Require().NoError(s.Run())
	ts.Require().DirExists(filepath.Join(s.plandir, plan.Manifest))
}

//nolint:paralleltest // cannot parallel because of setenv and uses helm repository.yaml flock
func TestNonParallelNonParallelBuildTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(NonParallelBuildTestSuite))
}
