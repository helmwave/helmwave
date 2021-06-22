package action

import (
	"fmt"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/tests"
	"github.com/urfave/cli/v2"
	"os"
	"testing"
)

func clean() {
	_ = os.RemoveAll(tests.Root + plan.Plandir)
}

func TestBuildRepositories(t *testing.T) {
	defer clean()

	s := &Build{
		plandir:  tests.Root + plan.Plandir,
		yml:      tests.Root + "02_helmwave.yml",
		tags:     cli.StringSlice{},
		matchAll: true,
	}

	err := s.Run()
	if err != nil {
		t.Error(err)
	}

	b, _ := plan.NewBody(tests.Root + plan.Plandir + plan.Planfile)

	if len(b.Repositories) != 1 && b.Repositories[0].Name != "bitnami" {
		t.Error("'bitnami' not found")
		fmt.Println(len(b.Repositories))
		fmt.Println(b.Repositories[0].Name)
	}

}

func TestBuildReleasesMatchGroup(t *testing.T) {
	defer clean()

	s := &Build{
		plandir:  tests.Root + plan.Plandir,
		yml:      tests.Root + "03_helmwave.yml",
		tags:     *cli.NewStringSlice("b"),
		matchAll: true,
	}

	err := s.Run()
	if err != nil {
		t.Error(err)
	}

	b, _ := plan.NewBody(tests.Root + plan.Plandir + plan.Planfile)

	if len(b.Releases) != 2 && b.Releases[0].Name != "redis-b" && b.Releases[1].Name != "memcached-b" {
		t.Error("'redis-b' and 'memcached-b' not found")
	}

}

func TestBuildReleasesMatchGroups(t *testing.T) {
	defer clean()

	s := &Build{
		plandir:  tests.Root + plan.Plandir,
		yml:      tests.Root + "03_helmwave.yml",
		tags:     *cli.NewStringSlice("b", "redis"),
		matchAll: true,
	}

	err := s.Run()
	if err != nil {
		t.Error(err)
	}

	b, _ := plan.NewBody(tests.Root + plan.Plandir + plan.Planfile)

	if len(b.Releases) != 1 && b.Releases[0].Name != "redis-b" {
		t.Error("'redis-b' not found")
	}

}
