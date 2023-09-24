package action

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/helmwave/helmwave/pkg/cache"
	"github.com/helmwave/helmwave/pkg/release"

	"github.com/helmwave/helmwave/pkg/plan"
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
}

func TestBuildTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildTestSuite))
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
		templater: template.TemplaterSprig,
	}

	s := &Build{
		yml:  y,
		tags: cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
	}
	createGenericFS(&s.yml.srcFS)
	createGenericFS(&s.yml.destFS, tests.Root, "helmwave.yml")
	createGenericFS(&s.planFS, tmpDir)

	ts.Require().Error(s.Run(context.Background()))
}

func (ts *BuildTestSuite) TestManifest() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		templater: template.TemplaterSprig,
	}

	s := &Build{
		yml:  y,
		tags: cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
	}
	createGenericFS(&s.yml.srcFS, tests.Root, "01_helmwave.yml.tpl")
	createGenericFS(&s.yml.destFS, tests.Root, "02_helmwave.yml")
	createGenericFS(&s.planFS, tmpDir)

	ts.Require().NoError(s.Run(context.Background()))
	ts.Require().DirExists(filepath.Join(tmpDir, plan.Manifest))
}

func (ts *BuildTestSuite) TestNonUniqueReleases() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		templater: template.TemplaterSprig,
	}
	createGenericFS(&y.srcFS, tests.Root, "14_helmwave.yml")
	createGenericFS(&y.destFS, tmpDir, "14_helmwave.yml")

	sfail := &Build{
		yml:  y,
		tags: cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml: true,
	}
	createGenericFS(&sfail.planFS, tmpDir)

	sfailByTag := &Build{
		yml:  y,
		tags: cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml: true,
	}
	createGenericFS(&sfailByTag.planFS, tmpDir)
	err := sfailByTag.tags.Set("nginx")
	ts.Require().NoError(err)

	sa := &Build{
		yml:  y,
		tags: cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml: true,
	}
	createGenericFS(&sa.planFS, tmpDir)
	err = sa.tags.Set("nginx-a")
	ts.Require().NoError(err)

	sb := &Build{
		yml:  y,
		tags: cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml: true,
	}
	createGenericFS(&sb.planFS, tmpDir)
	err = sb.tags.Set("nginx-b")
	ts.Require().NoError(err)

	ts.Require().ErrorIs(sfail.Run(context.Background()), release.DuplicateError{})
	ts.Require().ErrorIs(sfailByTag.Run(context.Background()), release.DuplicateError{})
	ts.Require().NoError(sa.Run(context.Background()))
	ts.Require().NoError(sb.Run(context.Background()))
}

func (ts *BuildTestSuite) TestRepositories() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		templater: template.TemplaterSprig,
	}

	s := &Build{
		yml:  y,
		tags: cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
	}
	createGenericFS(&s.yml.srcFS, tests.Root, "01_helmwave.yml.tpl")
	createGenericFS(&s.yml.destFS, tests.Root, "02_helmwave.yml")
	createGenericFS(&s.planFS, tmpDir)

	ts.Require().NoError(s.Run(context.Background()))

	const rep = "bitnami"
	planfileFS, _ := s.planFS.(fs.SubFS).Sub(plan.File)
	b, _ := plan.NewBody(context.Background(), planfileFS, true)

	if _, found := repo.IndexOfName(b.Repositories, rep); !found {
		ts.Failf("", "%q not found", rep)
	}
}

func (ts *BuildTestSuite) TestReleasesMatchGroup() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		templater: template.TemplaterSprig,
	}
	createGenericFS(&y.srcFS, tests.Root, "01_helmwave.yml.tpl")
	createGenericFS(&y.destFS, tests.Root, "03_helmwave.yml")

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
			yml:  y,
			tags: *cases[i].tags,
			options: plan.BuildOptions{
				MatchAll: true,
			},
		}
		createGenericFS(&s.planFS, tmpDir)

		ts.Require().NoError(s.Run(context.Background()))

		planfileFS, _ := s.planFS.(fs.SubFS).Sub(plan.File)
		b, _ := plan.NewBody(context.Background(), planfileFS, true)

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
		templater: template.TemplaterSprig,
	}

	s := &Build{
		tags: cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml:  true,
		yml:      y,
		diff:     &Diff{},
		diffMode: DiffModeLocal,
	}
	createGenericFS(&s.yml.srcFS, tests.Root, "07_helmwave.yml")
	createGenericFS(&s.yml.destFS, tests.Root, "07_helmwave.yml")
	createGenericFS(&s.planFS, tmpDir)

	ts.Require().NoError(s.Run(context.Background()), "build should not fail without diffing")
	ts.Require().NoError(s.Run(context.Background()), "build should not fail with diffing with previous plan")
}

type NonParallelBuildTestSuite struct {
	suite.Suite

	defaultHooks log.LevelHooks
	logHook      *logTest.Hook
}

//nolint:paralleltest // can't parallel because of setenv and uses helm repository.yaml flock
func TestNonParallelNonParallelBuildTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(NonParallelBuildTestSuite))
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
	res := make([]string, len(ts.logHook.Entries))
	for i, entry := range ts.logHook.AllEntries() {
		res[i] = entry.Message
	}

	return res
}

func (ts *NonParallelBuildTestSuite) TestAutoYml() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		templater: template.TemplaterSprig,
	}

	s := &Build{
		tags: cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml: true,
		yml:     y,
	}
	createGenericFS(&s.yml.srcFS, tests.Root, "01_helmwave.yml.tpl")
	createGenericFS(&s.yml.destFS, tmpDir, "01_auto_yaml_helmwave.yml")
	createGenericFS(&s.planFS, tmpDir)

	value := strings.ToLower(strings.ReplaceAll(ts.T().Name(), "/", ""))
	ts.T().Setenv("NAMESPACE", value)

	ts.Require().NoError(s.Run(context.Background()))
	ts.Require().DirExists(filepath.Join(tmpDir, plan.Manifest))
}

func (ts *NonParallelBuildTestSuite) TestGomplate() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		templater: template.TemplaterGomplate,
	}

	s := &Build{
		tags: cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml: true,
		yml:     y,
	}
	createGenericFS(&s.yml.srcFS, tests.Root, "08_helmwave.yml")
	createGenericFS(&s.yml.destFS, tmpDir, "08_helmwave.yml")
	createGenericFS(&s.planFS, tmpDir)

	ts.Require().NoError(s.Run(context.Background()))
	ts.Require().DirExists(filepath.Join(tmpDir, plan.Manifest))
}

func (ts *NonParallelBuildTestSuite) TestLifecycle() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		templater: template.TemplaterSprig,
	}

	s := &Build{
		tags: cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		autoYml: true,
		yml:     y,
	}
	createGenericFS(&s.yml.srcFS, tests.Root, "13_helmwave.yml")
	createGenericFS(&s.yml.destFS, tmpDir, "13_helmwave.yml")
	createGenericFS(&s.planFS, tmpDir)

	ts.Require().NoError(s.Run(context.Background()))
	ts.Require().DirExists(filepath.Join(tmpDir, plan.Manifest))

	logMessages := ts.getLoggerMessages()
	ts.Require().Contains(logMessages, "running pre_build script for nginx")
	ts.Require().Contains(logMessages, "run global pre_build script")
	ts.Require().Contains(logMessages, "running post_build script for nginx")
	ts.Require().Contains(logMessages, "run global post_build script")
}

func (ts *NonParallelBuildTestSuite) TestInvalidCacheDir() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		templater: template.TemplaterSprig,
	}

	s := &Build{
		yml:  y,
		tags: cli.StringSlice{},
		options: plan.BuildOptions{
			MatchAll: true,
		},
		chartsCacheDir: "/proc/1/bla",
	}
	createGenericFS(&s.yml.srcFS, tests.Root, "01_helmwave.yml.tpl")
	createGenericFS(&s.yml.destFS, tests.Root, "02_helmwave.yml")
	createGenericFS(&s.planFS, tmpDir)

	defer cache.ChartsCache.Init(nil, "") //nolint:errcheck
	ts.Require().Error(s.Run(context.Background()))
}
