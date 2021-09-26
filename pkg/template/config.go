package template

import (
	"github.com/hairyhenderson/gomplate/v3/data"
)

type Config struct {
	Gomplate GomplateConfig
}

type GomplateConfig struct {
	Data    *data.Data
	Enabled bool
}

var cfg *Config

func SetConfig(config *Config) {
	cfg = config
}
