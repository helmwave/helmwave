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

type SuccessMode string

const (
	SuccessModeIfEmpty  SuccessMode = "if_empty"
	SuccessModeIfVector SuccessMode = "if_vector"
)

type resultChecker func(result model.Value) error

// Config is the main monitor Config.
type Config struct {
	client  v1.API        `yaml:"-" json:"-"`
	log     *log.Entry    `yaml:"-" json:"-"`
	checker resultChecker `yaml:"-" json:"-"`

	URL         string      `yaml:"url" json:"url" jsonschema:"required,title=Prometheus URL"`
	Expr        string      `yaml:"expr" json:"expr" jsonschema:"required,title=Prometheus expression"`
	SuccessMode SuccessMode `yaml:"success_mode" json:"success_mode" jsonschema:"default=if_vector,title=Success mode,enum=if_empty,enum=if_vector"`
	Insecure    bool        `yaml:"insecure" json:"insecure" jsonschema:"default=false"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Init(_ context.Context, logger *log.Entry) error {
	client, err := api.NewClient(api.Config{Address: c.URL})
	if err != nil {
		return NewPrometheusClientError(err)
	}

	c.client = v1.NewAPI(client)
	c.log = logger

	if c.SuccessMode == "" {
		c.SuccessMode = SuccessModeIfVector
	}

	switch c.SuccessMode {
	case SuccessModeIfEmpty:
		c.checker = c.checkResultIfEmpty
	case SuccessModeIfVector, "":
		c.checker = c.checkResultIfVector
	default:
		return ErrInvalidSuccessMode
	}

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

	return c.checker(result)
}

//nolint:forcetypeassert
func (c *Config) checkResultIfEmpty(result model.Value) error {
	switch result.Type() {
	case model.ValNone:
		return nil
	case model.ValVector:
		v := result.(model.Vector)
		if len(v) == 0 {
			return nil
		}

		return ErrResultNotEmpty
	default:
		return ErrInvalidResult
	}
}

//nolint:forcetypeassert
func (c *Config) checkResultIfVector(result model.Value) error {
	switch result.Type() {
	case model.ValNone:
		return ErrResultEmpty
	case model.ValVector:
		v := result.(model.Vector)
		if len(v) == 0 {
			return ErrResultEmpty
		}

		return nil
	default:
		return ErrInvalidResult
	}
}

func (c *Config) Validate() error {
	if c.URL == "" {
		return ErrURLEmpty
	}

	if c.Expr == "" {
		return ErrExprEmpty
	}

	if c.SuccessMode != SuccessModeIfEmpty && c.SuccessMode != SuccessModeIfVector && c.SuccessMode != "" {
		return ErrInvalidSuccessMode
	}

	return nil
}
