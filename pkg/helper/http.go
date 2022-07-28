package helper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// HTTPTimeout is a timeout for HTTP requests.
	HTTPTimeout = 30 * time.Second
)

// Download downloads uri to file.
func Download(file, uri string) error {
	f, err := CreateFile(file)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Errorf("failed to close file %s: %v", f.Name(), err)
		}
	}(f)

	ctx, cancel := context.WithTimeout(context.Background(), HTTPTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", uri, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request to %s: %w", uri, err)
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", uri, err)
	}

	defer func(b io.Closer) {
		err := b.Close()
		if err != nil {
			log.Errorf("failed to close HTTP body: %v", err)
		}
	}(r.Body)

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
