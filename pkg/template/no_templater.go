package template

type noTemplater struct{}

func (t noTemplater) Name() string {
	return "copy"
}

func (t noTemplater) Render(src string, data any) ([]byte, error) {
	return []byte(src), nil
}

func (t noTemplater) Delims(string, string) {}
