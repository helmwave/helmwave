package action

import (
	"context"
	"github.com/helmwave/helmwave/pkg/templater"
	"path/filepath"
	"strings"
	"testing"

	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type StatusTestSuite struct {
	suite.Suite

	ctx context.Context
}

//nolint:paralleltest // can't parallel because of setenv
func TestStatusTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(StatusTestSuite))
}

func (ts *StatusTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *StatusTestSuite) TestCmd() {
	s := &Status{}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}

func (ts *StatusTestSuite) TestRun() {
	r := &Build{
		plandir: ts.T().TempDir(),
		tags:    cli.StringSlice{},
		autoYml: true,
		yml: &Yml{
			tpl:       filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
			file:      filepath.Join(ts.T().TempDir(), "02_helmwave.yml"),
			templater: templater.Default,
		},
	}

	value := strings.ToLower(strings.ReplaceAll(ts.T().Name(), "/", ""))
	ts.T().Setenv("NAMESPACE", value)

	ts.Require().NoError(r.Run(ts.ctx))

	s := &Status{
		build: r,
	}

	ts.Require().NoError(s.Run(ts.ctx))
}
