package action

import (
	"context"
	"strings"
	"testing"

	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type StatusTestSuite struct {
	suite.Suite
}

//nolint:paralleltest // can't parallel because of setenv
func TestStatusTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(StatusTestSuite))
}

func (ts *StatusTestSuite) TestCmd() {
	s := &Status{}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}

func (ts *StatusTestSuite) TestRun() {
	r := &Build{
		tags:    cli.StringSlice{},
		autoYml: true,
		yml: &Yml{
			templater: template.TemplaterSprig,
		},
	}
	createGenericFS(&r.yml.srcFS, tests.Root, "01_helmwave.yml.tpl")
	createGenericFS(&r.yml.destFS, ts.T().TempDir(), "02_helmwave.yml")
	createGenericFS(&r.planFS, ts.T().TempDir())
	createGenericFS(&r.contextFS, ts.T().TempDir())

	value := strings.ToLower(strings.ReplaceAll(ts.T().Name(), "/", ""))
	ts.T().Setenv("NAMESPACE", value)

	ts.Require().NoError(r.Run(context.Background()))

	s := &Status{
		build: r,
	}

	ts.Require().NoError(s.Run(context.Background()))
}
