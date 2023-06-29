package action

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type BuildTestSuite struct {
	suite.Suite
}

func (ts *BuildTestSuite) TestManifest() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
		file:      filepath.Join(tests.Root, "02_helmwave.yml"),
		templater: template.TemplaterSprig,
	}

	s := &Build{
		plandir: tmpDir,
		yml:     y,
		tags:    cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
	}

	ts.Require().NoError(s.Run(context.Background()))
	ts.Require().DirExists(filepath.Join(s.plandir, plan.Manifest))
}

// func (ts *BuildTestSuite) TestRepositories404() {
//	s := &Build{
//		plandir:  tmpDir,
//		ymlFile:      filepath.Join(tests.Root, "04_helmwave.yml"),
//		tags:     cli.StringSlice{},
//		options: plan.BuildOptions{
//				MatchAll: true,
//			},
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
		templater: template.TemplaterSprig,
	}

	s := &Build{
		plandir: tmpDir,
		yml:     y,
		tags:    cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
	}

	ts.Require().NoError(s.Run(context.Background()))

	const rep = "bitnami"
	b, _ := plan.NewBody(context.Background(), filepath.Join(s.plandir, plan.File))

	if _, found := repo.IndexOfName(b.Repositories, rep); !found {
		ts.Failf("%q not found", rep)
	}
}

func (ts *BuildTestSuite) TestReleasesMatchGroup() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
		file:      filepath.Join(tests.Root, "03_helmwave.yml"),
		templater: template.TemplaterSprig,
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
			plandir: tmpDir,
			yml:     y,
			tags:    *cases[i].tags,
			options: plan.BuildOptions{
				MatchAll: true,
			},
		}

		ts.Require().NoError(s.Run(context.Background()))

		b, _ := plan.NewBody(context.Background(), filepath.Join(s.plandir, plan.File))

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
		templater: template.TemplaterSprig,
	}

	s := &Build{
		plandir: tmpDir,
		tags:    cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml:  true,
		yml:      y,
		diff:     &Diff{},
		diffMode: DiffModeLocal,
	}

	ts.Require().NoError(s.Run(context.Background()), "build should not fail without diffing")
	ts.Require().NoError(s.Run(context.Background()), "build should not fail with diffing with previous plan")
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
		templater: template.TemplaterSprig,
	}

	s := &Build{
		plandir: tmpDir,
		tags:    cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml: true,
		yml:     y,
	}

	value := strings.ToLower(strings.ReplaceAll(ts.T().Name(), "/", ""))
	ts.T().Setenv("NAMESPACE", value)

	ts.Require().NoError(s.Run(context.Background()))
	ts.Require().DirExists(filepath.Join(s.plandir, plan.Manifest))
}

func (ts *NonParallelBuildTestSuite) TestGomplate() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "08_helmwave.yml"),
		file:      filepath.Join(tmpDir, "08_helmwave.yml"),
		templater: template.TemplaterGomplate,
	}

	s := &Build{
		plandir: tmpDir,
		tags:    cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml: true,
		yml:     y,
	}

	ts.Require().NoError(s.Run(context.Background()))
	ts.Require().DirExists(filepath.Join(s.plandir, plan.Manifest))
}

func (ts *NonParallelBuildTestSuite) TestLifecycle() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "13_helmwave.yml"),
		file:      filepath.Join(tmpDir, "13_helmwave.yml"),
		templater: template.TemplaterSprig,
	}

	s := &Build{
		plandir: tmpDir,
		tags:    cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml: true,
		yml:     y,
	}

	var buf bytes.Buffer
	oldOut := log.StandardLogger().Out
	log.StandardLogger().SetOutput(&buf)
	defer log.StandardLogger().SetOutput(oldOut)

	ts.Require().NoError(s.Run(context.Background()))
	ts.Require().DirExists(filepath.Join(s.plandir, plan.Manifest))

	output := buf.String()
	ts.Require().Contains(output, "running pre_build script for nginx")
	ts.Require().Contains(output, "run global pre_build script")
	ts.Require().Contains(output, "running post_build script for nginx")
	ts.Require().Contains(output, "run global post_build script")
}

//nolintlint:paralleltest // can't parallel because of setenv and uses helm repository.yaml flock
func TestNonParallelNonParallelBuildTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(NonParallelBuildTestSuite))
}
