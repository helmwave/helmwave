//go:build ignore || unit

package helper

import "testing"

func TestIsURL(t *testing.T) {
	urls := []string{
		"https://blog.golang.org/slices-intro",
		"https://helmwave.github.io/",
	}

	for _, url := range urls {
		b := IsURL(url)
		if !b {
			t.Error("bad url: ", url)
		}
	}
}
