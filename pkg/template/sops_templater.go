package template

import (
	"context"
	"io"

	"go.mozilla.org/sops/v3/decrypt"
)

const (
	TemplaterSOPS = "sops"
)

type sopsTemplater struct {
	additionalOutputs []io.Writer
}

func (t sopsTemplater) Name() string {
	return TemplaterSOPS
}

func (t sopsTemplater) Render(_ context.Context, src string, _ any) ([]byte, error) {
	data, err := decrypt.Data([]byte(src), "yaml")
	if err != nil {
		return nil, NewSOPSDecodeError(err)
	}

	w := io.MultiWriter(t.additionalOutputs...)

	_, err = w.Write(data)
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	return data, nil
}

func (t sopsTemplater) Delims(string, string) {}

func (t *sopsTemplater) AddOutput(w io.Writer) {
	if t.additionalOutputs == nil {
		t.additionalOutputs = []io.Writer{}
	}
	t.additionalOutputs = append(t.additionalOutputs, w)
}

func (t sopsTemplater) AddFunc(string, any) {}
