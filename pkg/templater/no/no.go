package no

import (
	"context"
	"io"
)

const (
	TemplaterName = "copy"
)

type Templater struct {
	additionalOutputs []io.Writer
}

func (t Templater) Name() string {
	return TemplaterName
}

func (t Templater) Render(_ context.Context, src []byte, _ any) ([]byte, error) {
	w := io.MultiWriter(t.additionalOutputs...)

	_, err := w.Write(src)
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	return src, nil
}

func (t Templater) Delims(string, string) {}

func (t *Templater) AddOutput(w io.Writer) {
	if t.additionalOutputs == nil {
		t.additionalOutputs = []io.Writer{}
	}
	t.additionalOutputs = append(t.additionalOutputs, w)
}

func (t Templater) AddFunc(string, any) {}
