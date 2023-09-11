package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"net/http"
	"slices"

	log "github.com/sirupsen/logrus"
)

const (
	TYPE          = "http"
	DefaultMethod = "HEAD"
)

// Config is the main monitor Config.
type Config struct {
	client        *http.Client      `yaml:"-" json:"-"`
	log           *log.Entry        `yaml:"-" json:"-"`
	URL           string            `yaml:"url" json:"url" jsonschema:"required,title=URL to query"`
	Method        string            `yaml:"method" json:"method" jsonschema:"title=HTTP method,default=HEAD"`
	Body          string            `yaml:"body" json:"body" jsonschema:"title=HTTP body,default="`
	Headers       map[string]string `yaml:"headers" json:"headers" jsonschema:"title=HTTP headers to set"`
	ExpectedCodes []int             `yaml:"expected_codes" json:"expected_codes" jsonschema:"required,title=Expected response codes"`
	Insecure      bool              `yaml:"insecure" json:"insecure" jsonschema:"default=false"`
}

func NewConfig() *Config {
	return &Config{
		Method: DefaultMethod,
	}
}

func (c *Config) Init(ctx context.Context, logger *log.Entry) error {
	c.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.Insecure,
			},
		},
		Timeout: 0,
	}

	c.log = logger

	return nil
}

func (c *Config) Run(ctx context.Context) error {
	body := bytes.NewBufferString(c.Body)

	req, err := http.NewRequestWithContext(ctx, c.Method, c.URL, body)
	if err != nil {
		return NewRequestError(err)
	}

	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return NewResponseError(err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			c.log.WithError(err).Error("failed to close HTTP body")
		}
	}()

	if !slices.Contains(c.ExpectedCodes, resp.StatusCode) {
		return NewUnexpectedStatusError(resp.StatusCode)
	}

	return nil
}

func (c *Config) Validate() error {
	if c.URL == "" {
		return ErrURLEmpty
	}

	return nil
}
