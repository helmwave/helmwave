package helper

import "testing"

func TestIsUrl(t *testing.T) {
	urls := []string{
		"https://blog.golang.org/slices-intro",
		"https://helmwave.github.io/",
	}

	for _, url := range urls {
		b := IsUrl(url)
		if !b {
			t.Error("bad url: ", url)
		}
	}
}
