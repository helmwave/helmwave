package action

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/databus23/helm-diff/v3/diff"
	"github.com/helmwave/helmwave/pkg/cache"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/hooks"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	log "github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type BuildTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestBuildTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildTestSuite))
}

func (ts *BuildTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *BuildTestSuite) TestCmd() {
	s := &Build{}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}

func (ts *BuildTestSuite) TestYmlError() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		file:      filepath.Join(tests.Root, "helmwave.yml"),
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

	ts.Require().Error(s.Run(ts.ctx))
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

	ts.Require().NoError(s.Run(ts.ctx))
	ts.Require().DirExists(filepath.Join(s.plandir, plan.Manifests))
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

func (ts *BuildTestSuite) TestNonUniqueReleases() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "14_helmwave.yml"),
		file:      filepath.Join(tmpDir, "14_helmwave.yml"),
		templater: template.TemplaterSprig,
	}

	sfail := &Build{
		plandir: tmpDir,
		yml:     y,
		tags:    cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml: true,
	}

	sfailByTag := &Build{
		plandir: tmpDir,
		yml:     y,
		tags:    cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml: true,
	}
	err := sfailByTag.tags.Set("nginx")
	ts.Require().NoError(err)

	sa := &Build{
		plandir: tmpDir,
		yml:     y,
		tags:    cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml: true,
	}
	err = sa.tags.Set("nginx-a")
	ts.Require().NoError(err)

	sb := &Build{
		plandir: tmpDir,
		yml:     y,
		tags:    cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml: true,
	}
	err = sb.tags.Set("nginx-b")
	ts.Require().NoError(err)

	var e *release.DuplicateError
	ts.Require().ErrorAs(sfail.Run(ts.ctx), &e)
	ts.Equal("nginx@test", e.Uniq.String())

	ts.Require().ErrorAs(sfailByTag.Run(ts.ctx), &e)
	ts.Equal("nginx@test", e.Uniq.String())

	ts.Require().NoError(sa.Run(ts.ctx))
	ts.Require().NoError(sb.Run(ts.ctx))
}

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

	ts.Require().NoError(s.Run(ts.ctx))

	const rep = "bitnami"
	b, _ := plan.NewBody(ts.ctx, filepath.Join(s.plandir, plan.File), true)

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

		ts.Require().NoError(s.Run(ts.ctx))

		b, _ := plan.NewBody(ts.ctx, filepath.Join(s.plandir, plan.File), true)

		names := helper.SlicesMap(b.Releases, func(rel release.Config) string {
			return rel.Name()
		})

		ts.ElementsMatch(cases[i].names, names)
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
		diff:     &Diff{Options: &diff.Options{}},
		diffMode: DiffModeLocal,
	}

	ts.Require().NoError(s.Run(ts.ctx), "build should not fail without diffing")
	ts.Require().NoError(s.Run(ts.ctx), "build should not fail with diffing with previous plan")
}

func (ts *BuildTestSuite) TestValuesDependency() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "19_helmwave.yml"),
		file:      filepath.Join(tests.Root, "19_helmwave.yml"),
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

	ts.Require().NoError(s.Run(ts.ctx), "build should not fail")
}

type NonParallelBuildTestSuite struct {
	suite.Suite

	defaultHooks log.LevelHooks
	logHook      *logTest.Hook

	ctx context.Context
}

//nolint:paralleltest // can't parallel because of setenv and uses helm repository.yaml flock
func TestNonParallelNonParallelBuildTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(NonParallelBuildTestSuite))
}

func (ts *NonParallelBuildTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *NonParallelBuildTestSuite) SetupSuite() {
	ts.defaultHooks = log.StandardLogger().Hooks
	ts.logHook = logTest.NewLocal(log.StandardLogger())
}

func (ts *NonParallelBuildTestSuite) TearDownTestSuite() {
	ts.logHook.Reset()
}

func (ts *NonParallelBuildTestSuite) TearDownSuite() {
	log.StandardLogger().ReplaceHooks(ts.defaultHooks)
}

func (ts *NonParallelBuildTestSuite) getLoggerMessages() []string {
	return helper.SlicesMap(ts.logHook.AllEntries(), func(entry *log.Entry) string {
		return entry.Message
	})
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

	ts.Require().NoError(s.Run(ts.ctx))
	ts.Require().DirExists(filepath.Join(s.plandir, plan.Manifests))
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

	ts.Require().NoError(s.Run(ts.ctx))
	ts.Require().DirExists(filepath.Join(s.plandir, plan.Manifests))
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

	ts.Require().NoError(s.Run(ts.ctx))
	ts.Require().DirExists(filepath.Join(s.plandir, plan.Manifests))

	logMessages := ts.getLoggerMessages()
	ts.Require().Contains(logMessages, "running pre_build script for nginx")
	ts.Require().Contains(logMessages, "run global pre_build script")
	ts.Require().Contains(logMessages, "running post_build script for nginx")
	ts.Require().Contains(logMessages, "run global post_build script")
}

func (ts *NonParallelBuildTestSuite) TestLifecyclePost() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "17_helmwave.yml"),
		file:      filepath.Join(tmpDir, "17_helmwave.yml"),
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

	err := s.Run(ts.ctx)

	var e *hooks.CommandRunError
	ts.Require().ErrorAs(err, &e)
}

func (ts *NonParallelBuildTestSuite) TestRemoteSource() {
	tmpDir := ts.T().TempDir()
	d, err := os.Getwd()
	ts.Require().NoError(err)

	ts.Require().NoError(os.Chdir(tmpDir))
	ts.T().Cleanup(func() {
		ts.Require().NoError(os.Chdir(d))
	})

	y := &Yml{
		tpl:       filepath.Join("tests", "02_helmwave.yml"),
		file:      filepath.Join("tests", "02_helmwave.yml"),
		templater: template.TemplaterSprig,
	}

	s := &Build{
		plandir: plan.Dir,
		tags:    cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		remoteSource: "github.com/helmwave/helmwave",
		yml:          y,
	}

	cache.DefaultConfig.Home = ts.T().TempDir()

	err = s.Run(ts.ctx)
	ts.Require().NoError(err)

	ts.DirExists(filepath.Join(tmpDir, s.plandir))
	ts.FileExists(filepath.Join(tmpDir, s.plandir, plan.File))
}
