package helper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Download(file, uri string) error {
	f, err := CreateFile(file)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.TODO(), "GET", uri, nil)
	if err != nil {
		return err
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", r.Status)
	}

	_, err = io.Copy(f, r.Body)
	if err != nil {
		return err
	}

	return nil
}

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
