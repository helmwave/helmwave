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
