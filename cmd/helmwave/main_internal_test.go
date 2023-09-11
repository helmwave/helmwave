package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type CliTestSuite struct {
	suite.Suite
}

func TestCliTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CliTestSuite))
}

//nolint:unparam // we may use not all buffers
func (ts *CliTestSuite) prepareApp() (*cli.App, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	ts.T().Helper()

	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := CreateApp()
	app.Reader = stdin
	app.Writer = stdout
	app.ErrWriter = stderr

	return app, stdin, stdout, stderr
}

func (ts *CliTestSuite) TestCommandNotFound() {
	app, _, _, _ := ts.prepareApp() //nolint:dogsled // no need to access nor stdin or stdout or stderr

	cmd := ts.T().Name()
	expectedError := CommandNotFoundError{Command: cmd}.Error()

	ts.Require().PanicsWithError(
		expectedError,
		func() {
			_ = app.Run([]string{"helmwave", cmd})
		},
	)
}

func (ts *CliTestSuite) TestCommandsList() {
	requiredCommands := []string{"build", "up", "down", "yml"}

	app, _, _, _ := ts.prepareApp() //nolint:dogsled // no need to access nor stdin or stdout or stderr

	commands := app.VisibleCommands()
	cmds := make([]string, 0, len(commands))

	for _, cmd := range commands {
		cmds = append(cmds, cmd.Name)
		cmds = append(cmds, cmd.Aliases...)
	}

	ts.Require().Subset(cmds, requiredCommands)
}

func (ts *CliTestSuite) TestRecoverWithoutPanic() {
	ts.Require().NotPanics(recoverPanic)
}

func (ts *CliTestSuite) TestRecoverPanic() {
	err := errors.New(ts.T().Name())
	ts.Require().Panics(func() {
		defer recoverPanic()
		panic(err)
	})
}
