package template

import (
	"go.mozilla.org/sops/v3/decrypt"
)

const (
	TemplaterSOPS = "sops"
)

type sopsTemplater struct{}

func (t sopsTemplater) Name() string {
	return TemplaterSOPS
}

func (t sopsTemplater) Render(src string, _ any) ([]byte, error) {
	data, err := decrypt.Data([]byte(src), "yaml")
	if err != nil {
		return nil, NewSOPSDecodeError(err)
	}

	return data, nil
}

func (t sopsTemplater) Delims(string, string) {}
