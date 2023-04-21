package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type CliTestSuite struct {
	suite.Suite
}

// nolintlint:unparam // we may use not all buffers
func (s *CliTestSuite) prepareApp() (*cli.App, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	s.T().Helper()

	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := CreateApp()
	app.Reader = stdin
	app.Writer = stdout
	app.ErrWriter = stderr

	return app, stdin, stdout, stderr
}

func (s *CliTestSuite) TestCommandNotFound() {
	app, _, _, _ := s.prepareApp() // nolintlint:dogsled // no need to access nor stdin or stdout or stderr

	cmd := s.T().Name()
	expectedError := CommandNotFoundError{Command: cmd}.Error()

	s.Require().PanicsWithError(
		expectedError,
		func() {
			_ = app.Run([]string{"helmwave", cmd})
		},
	)
}

func (s *CliTestSuite) TestCommandsList() {
	requiredCommands := []string{"build", "up", "down", "yml"}

	app, _, _, _ := s.prepareApp() // nolintlint:dogsled // no need to access nor stdin or stdout or stderr

	commands := app.VisibleCommands()
	cmds := make([]string, 0, len(commands))

	for _, cmd := range commands {
		cmds = append(cmds, cmd.Name)
		cmds = append(cmds, cmd.Aliases...)
	}

	s.Require().Subset(cmds, requiredCommands)
}

func TestCliTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CliTestSuite))
}
