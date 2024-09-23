package sops

import (
	"context"
	"io"

	"github.com/getsops/sops/v3/decrypt"
)

const (
	TemplaterName = "sops"
)

type Templater struct {
	additionalOutputs []io.Writer
}

func (t Templater) Name() string {
	return TemplaterName
}

func (t Templater) Render(_ context.Context, src []byte, _ any) ([]byte, error) {
	data, err := decrypt.Data(src, "yaml")
	if err != nil {
		return nil, NewDecodeError(err)
	}

	w := io.MultiWriter(t.additionalOutputs...)

	_, err = w.Write(data)
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	return data, nil
}

func (t Templater) Delims(string, string) {}

func (t *Templater) AddOutput(w io.Writer) {
	if t.additionalOutputs == nil {
		t.additionalOutputs = []io.Writer{}
	}
	t.additionalOutputs = append(t.additionalOutputs, w)
}

func (t Templater) AddFunc(string, any) {}
