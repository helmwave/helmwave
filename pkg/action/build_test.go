package action

import (
	"os"
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/tests"
	"github.com/urfave/cli/v2"
)

func clean() {
	_ = os.RemoveAll(tests.Root + plan.Dir)
}

func TestBuildManifest(t *testing.T) {
	defer clean()

	s := &Build{
		plandir:  tests.Root + plan.Dir,
		yml:      tests.Root + "02_helmwave.yml",
		tags:     cli.StringSlice{},
		matchAll: true,
	}

	err := s.Run()
	if err != nil {
		t.Error(err)
	}

	if ok := helper.IsExists(tests.Root + plan.Dir + plan.Manifest); !ok {
		t.Error(plan.ErrManifestDirNotFound)
	}
}

func TestBuildRepositories404(t *testing.T) {
	defer clean()

	s := &Build{
		plandir:  tests.Root + plan.Dir,
		yml:      tests.Root + "04_helmwave.yml",
		tags:     cli.StringSlice{},
		matchAll: true,
	}

	err := s.Run()
	if err != repo.ErrNotFound && err != nil {
		t.Error("'bitnami' must be not found")
	}
}

func TestBuildRepositories(t *testing.T) {
	defer clean()

	s := &Build{
		plandir:  tests.Root + plan.Dir,
		yml:      tests.Root + "02_helmwave.yml",
		tags:     cli.StringSlice{},
		matchAll: true,
	}

	err := s.Run()
	if err != nil {
		t.Error(err)
	}

	b, _ := plan.NewBody(tests.Root + plan.Dir + plan.File)

	if len(b.Repositories) != 1 && b.Repositories[0].Name != "bitnami" {
		t.Error("'bitnami' not found")
	}
}

func TestBuildReleasesMatchGroup(t *testing.T) {
	defer clean()

	s := &Build{
		plandir:  tests.Root + plan.Dir,
		yml:      tests.Root + "03_helmwave.yml",
		tags:     *cli.NewStringSlice("b"),
		matchAll: true,
	}

	err := s.Run()
	if err != nil {
		t.Error(err)
	}

	b, _ := plan.NewBody(tests.Root + plan.Dir + plan.File)

	if len(b.Releases) != 2 && b.Releases[0].Name != "redis-b" && b.Releases[1].Name != "memcached-b" {
		t.Error("'redis-b' and 'memcached-b' not found")
	}
}

func TestBuildReleasesMatchGroups(t *testing.T) {
	defer clean()

	s := &Build{
		plandir:  tests.Root + plan.Dir,
		yml:      tests.Root + "03_helmwave.yml",
		tags:     *cli.NewStringSlice("b", "redis"),
		matchAll: true,
	}

	err := s.Run()
	if err != nil {
		t.Error(err)
	}

	b, _ := plan.NewBody(tests.Root + plan.Dir + plan.File)

	if len(b.Releases) != 1 && b.Releases[0].Name != "redis-b" {
		t.Error("'redis-b' not found")
	}
}
