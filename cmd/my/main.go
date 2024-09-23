package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hairyhenderson/gomplate/v4"
	"net/url"
)

func main() {
	ctx := context.Background()

	content, result := &bytes.Buffer{}, &bytes.Buffer{}

	u, _ := url.Parse("https://raw.githubusercontent.com/helmwave/helmwave/refs/heads/main/tests/08_values.yaml")

	tr := gomplate.NewRenderer(gomplate.RenderOptions{
		Context: map[string]gomplate.DataSource{
			"content": {URL: u},
		},
	})

	tr.Render(ctx, "content", `{{ .content }}`, content)
	tr.Render(ctx, "rendered", content.String(), result)

	fmt.Println(result)

}

func isURI(uri string) bool {
	_, err := url.ParseRequestURI(uri)
	return err == nil
}
