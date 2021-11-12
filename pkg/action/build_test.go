//go:build ignore || integration

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

	err := s.Run()
	if err != nil {
		t.Error(err)
	}

	if ok := helper.IsExists(tests.Root + plan.Dir + plan.Manifest); !ok {
		t.Error(plan.ErrManifestDirNotFound)
	}
}

// func TestBuildRepositories404(t *testing.T) {
//	defer clean()
//
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

func TestBuildRepositories(t *testing.T) {
	defer clean()

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

	err := s.Run()
	if err != nil {
		t.Error(err)
	}

	const rep = "bitnami"
	b, _ := plan.NewBody(tests.Root + plan.Dir + plan.File)

	if _, found := repo.IndexOfName(b.Repositories, rep); !found {
		t.Errorf("%q not found", rep)
	}
}

func TestBuildReleasesMatchGroup(t *testing.T) {
	defer clean()

	y := &Yml{
		tests.Root + "01_helmwave.yml.tpl",
		tests.Root + "03_helmwave.yml",
	}

	s := &Build{
		plandir:  tests.Root + plan.Dir,
		yml:      y,
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

	y := &Yml{
		tests.Root + "01_helmwave.yml.tpl",
		tests.Root + "03_helmwave.yml",
	}

	s := &Build{
		plandir:  tests.Root + plan.Dir,
		yml:      y,
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

func TestBuildAutoYml(t *testing.T) {
	defer clean()

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

	err := s.Run()
	if err != nil {
		t.Error(err)
	}

	if ok := helper.IsExists(tests.Root + plan.Dir + plan.Manifest); !ok {
		t.Error(plan.ErrManifestDirNotFound)
	}
}

func TestBuildGomplate(t *testing.T) {
	defer clean()

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

	err := s.Run()
	if err != nil {
		t.Error(err)
	}

	if ok := helper.IsExists(tests.Root + plan.Dir + plan.Manifest); !ok {
		t.Error(plan.ErrManifestDirNotFound)
	}
}
