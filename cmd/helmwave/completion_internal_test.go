package main

func (ts *CliTestSuite) TestCompletion() {
	tests := []struct {
		args  []string
		fails bool
	}{
		{
			args:  []string{"helmwave", "completion"},
			fails: true,
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
			args:  []string{"helmwave", "completion", "ash"},
			fails: true,
		},
	}

	app, _, _, _ := ts.prepareApp() //nolint:dogsled // no need to access nor stdin or stdout or stderr

	for i := range tests {
		tt := tests[i]

		err := app.Run(tt.args)
		if tt.fails {
			ts.Error(err)
		} else {
			ts.NoError(err)
		}
	}
}
