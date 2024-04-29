package template

import (
	"context"
	"io"
)

const (
	TemplaterNone = "copy"
)

type noTemplater struct {
	additionalOutputs []io.Writer
}

func (t noTemplater) Name() string {
	return TemplaterNone
}

func (t noTemplater) Render(_ context.Context, src string, _ any) ([]byte, error) {
	w := io.MultiWriter(t.additionalOutputs...)

	b := []byte(src)
	_, err := w.Write(b)
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	return b, nil
}

func (t noTemplater) Delims(string, string) {}

func (t *noTemplater) AddOutput(w io.Writer) {
	if t.additionalOutputs == nil {
		t.additionalOutputs = []io.Writer{}
	}
	t.additionalOutputs = append(t.additionalOutputs, w)
}

func (t noTemplater) AddFunc(string, any) {}
