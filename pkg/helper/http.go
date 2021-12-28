package helper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Download downloads uri to file.
func Download(file, uri string) error {
	f, err := CreateFile(file)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.TODO(), "GET", uri, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request to %s: %w", uri, err)
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", uri, err)
	}

	defer r.Body.Close() //nolint:errcheck // TODO: need to check error

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", r.Status)
	}

	_, err = io.Copy(f, r.Body)
	if err != nil {
		return fmt.Errorf("failed to copy body of %s: %w", uri, err)
	}

	return nil
}

// IsURL reports whether provided string is a valid URL.
func IsURL(str string) bool {
	u, err := url.Parse(str)

	return err == nil && u.Scheme != "" && u.Host != ""
}
