package template

const (
	TemplaterNone = "copy"
)

type noTemplater struct{}

func (t noTemplater) Name() string {
	return TemplaterNone
}

func (t noTemplater) Render(src string, data any) ([]byte, error) {
	return []byte(src), nil
}

func (t noTemplater) Delims(string, string) {}
