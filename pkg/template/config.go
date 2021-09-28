package template

import (
	"net/url"

	"github.com/hairyhenderson/gomplate/v3/data"
)

type Config struct {
	Gomplate GomplateConfig
}

type GomplateConfig struct {
	Datasources map[string]Source
	data        *data.Data
	Enabled     bool
}

type Source struct {
	URL *url.URL
}

func (s *Source) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	type raw struct {
		URL string
	}

	r := raw{}
	err := unmarshal(&r)
	if err != nil {
		return err
	}

	u, err := url.Parse(r.URL)
	if err != nil {
		return err
	}

	*s = Source{
		URL: u,
	}
	return nil
}

var cfg *Config

func SetConfig(config *Config) {
	cfg = config

	if cfg == nil {
		return
	}

	sources := map[string]*data.Source{}
	for k, v := range cfg.Gomplate.Datasources {
		sources[k] = &data.Source{
			Alias: k,
			URL:   v.URL,
		}
	}

	cfg.Gomplate.data = &data.Data{
		Sources: sources,
	}
}
