package template

import "context"

const (
	TemplaterNone = "copy"
)

type noTemplater struct{}

func (t noTemplater) Name() string {
	return TemplaterNone
}

func (t noTemplater) Render(_ context.Context, src string, _ any) ([]byte, error) {
	return []byte(src), nil
}

func (t noTemplater) Delims(string, string) {}
