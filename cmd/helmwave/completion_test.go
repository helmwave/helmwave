package main_test

func (s *CliTestSuite) TestCompletion() {
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

	app, _, _, _ := s.prepareApp() //nolint:dogsled // no need to access nor stdin or stdout or stderr

	for _, tt := range tests {
		err := app.Run(tt.args)
		if tt.fails {
			s.Error(err)
		} else {
			s.NoError(err)
		}
	}
}
