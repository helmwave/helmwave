package main

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
			ts.Run("fails case", func() {
				ts.Assert().Panics(func() {
					err := app.Run(tt.args)
					ts.Assert().NoError(err, "unexpected error occurred")
				}, "Expected panic when args are: %v", tt.args)
			})
		} else {
			ts.Run("success case", func() {
				ts.Assert().NotPanics(func() {
					err := app.Run(tt.args)
					ts.Assert().NoError(err, "unexpected error occurred")
				}, "Unexpected panic for args: %v", tt.args)
			})
		}
	}
}
