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

	app, _, _, _ := ts.prepareApp() //nolint:dogsled // no need to access nor stdin/stderr

	// Avoid copying structs by using indices.
	for i := range tests {
		tt := &tests[i] // Take a pointer to the struct instead of copying it.

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
