package action

import (
	"context"
	"path/filepath"
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

//nolintlint:paralleltest // can't parallel because of setenv
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
		plandir: ts.T().TempDir(),
		tags:    cli.StringSlice{},
		autoYml: true,
		yml: &Yml{
			tpl:       filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
			file:      filepath.Join(ts.T().TempDir(), "02_helmwave.yml"),
			templater: template.TemplaterSprig,
		},
	}

	value := strings.ToLower(strings.ReplaceAll(ts.T().Name(), "/", ""))
	ts.T().Setenv("NAMESPACE", value)

	ts.Require().NoError(r.Run(context.Background()))

	s := &Status{
		build: r,
	}

	ts.Require().NoError(s.Run(context.Background()))
}
