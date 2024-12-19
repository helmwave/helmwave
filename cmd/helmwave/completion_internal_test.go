package main

import "github.com/stretchr/testify/assert"

func (ts *CliTestSuite) TestCompletion() {
	tests := []struct {
		args  []string
		fails bool
	}{
		{
			args:  []string{"helmwave", "completion"},
			fails: false,
		},
		{
			args:  []string{"helmwave", "completion", "bash"},
			fails: false,
		},
		{
			args:  []string{"helmwave", "completion", "zsh"},
			fails: false,
		},
		{
			args:  []string{"helmwave", "completion", "fish"},
			fails: false,
		},
		{
			args:  []string{"helmwave", "completion", "ash"},
			fails: true,
		},
	}

	app, _, _, _ := ts.prepareApp() //nolint:dogsled // no need to access nor stdin or stdout or stderr

	for _, tt := range tests {
		if tt.fails {
			assert.Panics(ts.T(), func() {
				err := app.Run(tt.args)
				assert.NoError(ts.T(), err, "unexpected error occurred")
			}, "Expected panic when args are: %v", tt.args)
		} else {
			assert.NotPanics(ts.T(), func() {
				err := app.Run(tt.args)
				assert.NoError(ts.T(), err, "unexpected error occurred")
			}, "Unexpected panic for args: %v", tt.args)
		}
	}
}
