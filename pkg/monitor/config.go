package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/helmwave/helmwave/pkg/monitor/prometheus"
	"github.com/invopop/jsonschema"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	DefaultTimeout = time.Minute * 5
)

// Config is the main monitor Config.
type config struct {
	Prometheus *prometheus.Config `yaml:",inline" json:",inline"`
	subConfig  SubConfig          `yaml:"-" json:"-"`
	log        *log.Entry         `yaml:"-" json:"-"`
	NameF      string             `yaml:"name" json:"name" jsonschema:"required"`
	Type       string             `yaml:"type" json:"type" jsonschema:"enum=prometheus"`
	Timeout    time.Duration      `yaml:"timeout" json:"timeout" jsonschema:"default=5m"`
}

type typeConfig struct {
	Type string `yaml:"type" json:"type"`
}

type _config config

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
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	c.Logger().Debug("started monitor")

	return c.subConfig.Run(ctx, c.Logger())
}

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse subconfig.
func (r *config) UnmarshalYAML(node *yaml.Node) error {
	t := typeConfig{}
	err := node.Decode(&t)
	if err != nil {
		return NewYAMLDecodeError(err)
	}

	r.Timeout = DefaultTimeout
	cfg := (*_config)(r)

	switch t.Type {
	case prometheus.TYPE:
		cfg.Prometheus = prometheus.NewConfig()
		cfg.subConfig = cfg.Prometheus
	default:
		return fmt.Errorf("unknown monitor type %q", t.Type)
	}

	err = node.Decode(cfg)
	if err != nil {
		return NewYAMLDecodeError(err)
	}

	return nil
}

//nolint:forcetypeassert
func (config) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
	}

	schema := r.Reflect(&_config{})

	pSchema, _ := schema.Properties.Get("Prometheus")
	promSchema := pSchema.(*jsonschema.Schema)
	schema.Properties.Delete("Prometheus")
	for _, k := range promSchema.Properties.Keys() {
		v, _ := promSchema.Properties.Get(k)
		schema.Properties.Set(k, v)
	}

	return schema
}
