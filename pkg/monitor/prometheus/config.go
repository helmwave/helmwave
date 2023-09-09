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
	TYPE                      = "prometheus"
	DEFAULT_TIMEOUT           = time.Minute
	DEFAULT_SUCCESS_THRESHOLD = 3
	DEFAULT_FAILURE_THRESHOLD = 3
)

// Config is the main monitor Config.
type Config struct {
	URL              string        `yaml:"url" json:"url" jsonschema:"required,title=Prometheus URL"`
	Expr             string        `yaml:"expr" json:"expr" jsonschema:"required,title=Prometheus expression"`
	Interval         time.Duration `yaml:"interval" json:"interval" jsonschema:"default=1m"`
	SuccessThreshold uint8         `yaml:"success_threshold" json:"success_threshold" jsonschema:"default=3"`
	FailureThreshold uint8         `yaml:"failure_threshold" json:"failure_threshold" jsonschema:"default=3"`
	Insecure         bool          `yaml:"insecure" json:"insecure" jsonschema:"default=false"`
}

func NewConfig() *Config {
	return &Config{
		Interval:         DEFAULT_TIMEOUT,
		SuccessThreshold: DEFAULT_SUCCESS_THRESHOLD,
		FailureThreshold: DEFAULT_FAILURE_THRESHOLD,
	}
}

func (c *Config) getAPIClient() (v1.API, error) {
	client, err := api.NewClient(api.Config{Address: c.URL})
	if err != nil {
		return nil, NewPrometheusClientError(err)
	}

	return v1.NewAPI(client), nil
}

func (c *Config) Run(ctx context.Context, logger *log.Entry) error {
	client, err := c.getAPIClient()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()

	var successStreak uint8 = 0
	var failureStreak uint8 = 0

	for (successStreak < c.SuccessThreshold) && (failureStreak < c.FailureThreshold) {
		select {
		case <-ticker.C:
			result, l, succeeded := c.runQuery(ctx, logger, client)

			if succeeded {
				successStreak += 1
				failureStreak = 0
				l.WithField("streak", successStreak).Debug("monitor succeeded")
			} else {
				successStreak = 0
				failureStreak += 1
				l.WithField("streak", failureStreak).WithField("result", result).Debug("monitor did not succeed")
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if failureStreak > 0 {
		return ErrFailureStreak
	}

	return nil
}

func (c *Config) runQuery(ctx context.Context, logger *log.Entry, client v1.API) (model.Value, *log.Entry, bool) {
	result, warns, err := client.Query(ctx, c.Expr, time.Now(), v1.WithTimeout(c.Interval))
	l := logger
	succeeded := true

	if err != nil {
		l = l.WithError(err)
		succeeded = false
	}

	if len(warns) > 0 {
		l = l.WithField("warnings", warns)
	}

	logger.WithField("result", result).Trace("monitor response")

	v, ok := result.(model.Vector)
	if !ok {
		l.Warn("failed to get result as vector")
		succeeded = false
	}

	if len(v) == 0 {
		succeeded = false
	}

	return result, l, succeeded
}
