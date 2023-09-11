package prometheus

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
)

const (
	TYPE = "prometheus"
)

// Config is the main monitor Config.
type Config struct {
	client   v1.API     `yaml:"-" json:"-"`
	log      *log.Entry `yaml:"-" json:"-"`
	URL      string     `yaml:"url" json:"url" jsonschema:"required,title=Prometheus URL"`
	Expr     string     `yaml:"expr" json:"expr" jsonschema:"required,title=Prometheus expression"`
	Insecure bool       `yaml:"insecure" json:"insecure" jsonschema:"default=false"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Init(ctx context.Context, logger *log.Entry) error {
	client, err := api.NewClient(api.Config{Address: c.URL})
	if err != nil {
		return NewPrometheusClientError(err)
	}

	c.client = v1.NewAPI(client)
	c.log = logger

	return nil
}

func (c *Config) Run(ctx context.Context) error {
	l := c.log

	now := time.Now()

	result, warns, err := c.client.Query(ctx, c.Expr, now)

	if len(warns) > 0 {
		l = l.WithField("warnings", warns)
	}

	if err != nil {
		return NewPrometheusClientError(err)
	}

	l.WithField("result", result).Trace("monitor response")

	v, ok := result.(model.Vector)
	if !ok {
		err = ErrResultNotVector
	}

	if len(v) == 0 {
		err = ErrResultEmpty
	}

	return err
}

func (c *Config) Validate() error {
	if c.URL == "" {
		return ErrURLEmpty
	}

	if c.Expr == "" {
		return ErrExprEmpty
	}

	return nil
}
