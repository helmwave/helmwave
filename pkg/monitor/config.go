package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/helmwave/helmwave/pkg/monitor/http"
	"github.com/helmwave/helmwave/pkg/monitor/prometheus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	DefaultTotalTimeout     = time.Minute * 5
	DefaultIterationTimeout = time.Second * 10
	DefaultInterval         = time.Minute
	DefaultSuccessThreshold = 3
	DefaultFailureThreshold = 3
)

// Config is the main monitor Config.
type config struct {
	Prometheus       *prometheus.Config `yaml:"prometheus" json:"prometheus" jsonschema:"title=Config for prometheus type,oneof_required=prometheus"`
	HTTP             *http.Config       `yaml:"http" json:"http" jsonschema:"title=Config for http type,oneof_required=http"`
	subConfig        SubConfig          `yaml:"-" json:"-"`
	log              *log.Entry         `yaml:"-" json:"-"`
	NameF            string             `yaml:"name" json:"name" jsonschema:"required"`
	Type             string             `yaml:"type" json:"type" jsonschema:"enum=prometheus,enum=http,required"`
	TotalTimeout     time.Duration      `yaml:"total_timeout" json:"total_timeout" jsonschema:"title=Timeout for the whole monitor,description=After this timeout hits monitor will fail regardless of current streak,default=5m"`
	IterationTimeout time.Duration      `yaml:"iteration_timeout" json:"iteration_timeout" jsonschema:"title=Timeout for each timeout execution,description=After this timeout hits monitor iteration will be considered as failed,default=10s"`
	Interval         time.Duration      `yaml:"interval" json:"interval" jsonschema:"default=1m"`
	SuccessThreshold uint8              `yaml:"success_threshold" json:"success_threshold" jsonschema:"default=3"`
	FailureThreshold uint8              `yaml:"failure_threshold" json:"failure_threshold" jsonschema:"default=3"`
}

type typeConfig struct {
	Type string `yaml:"type" json:"type"`
}

type _config config

func (c *config) setDefaults() {
	c.TotalTimeout = DefaultTotalTimeout
	c.IterationTimeout = DefaultIterationTimeout
	c.Interval = DefaultInterval
	c.SuccessThreshold = DefaultSuccessThreshold
	c.FailureThreshold = DefaultFailureThreshold
}

func (c *config) Name() string {
	return c.NameF
}

func (c *config) Logger() *log.Entry {
	if c.log == nil {
		c.log = log.WithField("monitor", c.Name())
	}

	return c.log
}

func (c *config) Run(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.TotalTimeout)
	defer cancel()

	c.Logger().Debug("initializing monitor")
	err := c.subConfig.Init(ctx, c.Logger())
	if err != nil {
		return NewMonitorInitError(err)
	}

	c.Logger().Debug("starting monitor")

	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()

	var successStreak uint8 = 0
	var failureStreak uint8 = 0

	for (successStreak < c.SuccessThreshold) && (failureStreak < c.FailureThreshold) {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(ctx, c.IterationTimeout)
			err := c.subConfig.Run(ctx)
			cancel()

			if err == nil {
				successStreak += 1
				failureStreak = 0
				c.Logger().
					WithField("streak", fmt.Sprintf("%d/%d", successStreak, c.SuccessThreshold)).
					Info("monitor succeeded")
			} else {
				successStreak = 0
				failureStreak += 1
				c.Logger().
					WithField("streak", fmt.Sprintf("%d/%d", failureStreak, c.FailureThreshold)).
					WithField("error", err).
					Info("monitor did not succeed")
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

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse subconfig.
func (c *config) UnmarshalYAML(node *yaml.Node) error {
	t := typeConfig{}
	err := node.Decode(&t)
	if err != nil {
		return NewYAMLDecodeError(err)
	}

	c.setDefaults()

	cfg := (*_config)(c)

	switch t.Type {
	case prometheus.TYPE:
		cfg.Prometheus = prometheus.NewConfig()
		cfg.subConfig = cfg.Prometheus
	case http.TYPE:
		cfg.HTTP = http.NewConfig()
		cfg.subConfig = cfg.HTTP
	default:
		return fmt.Errorf("unknown monitor type %q", t.Type)
	}

	err = node.Decode(cfg)
	if err != nil {
		return NewYAMLDecodeError(err)
	}

	return nil
}
