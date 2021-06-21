package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/tests"
	"github.com/urfave/cli/v2"
	"testing"
)

func TestBuild(t *testing.T) {
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

}
