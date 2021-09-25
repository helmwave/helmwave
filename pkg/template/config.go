package template

import (
	"github.com/hairyhenderson/gomplate/v3/data"
)

type Config struct {
	Gomplate GomplateConfig
}

type GomplateConfig struct {
	Enabled bool

	Data *data.Data
}

var (
	cfg *Config
)

func SetConfig(config *Config) {
	cfg = config
}
