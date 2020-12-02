package helmwave

import (
	log "github.com/sirupsen/logrus"
	"github.com/zhilyaev/helmwave/pkg/formatter"
)

type Log struct {
	//Engine *log.Logger
	Level  string
	Format string
	Color  bool
}

func (c *Config) InitLogger() error {
	c.InitLoggerFormat()
	return c.InitLoggerLevel()
}

func (c *Config) InitLoggerLevel() error {
	level, err := log.ParseLevel(c.Logger.Level)
	if err != nil {
		return err
	}
	//c.Logger.Engine.SetLevel(level)
	log.SetLevel(level)

	return nil
}

func (c *Config) InitLoggerFormat() {
	switch c.Logger.Format {
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			PrettyPrint: true,
		})
	case "pad":
		log.SetFormatter(&log.TextFormatter{
			PadLevelText: true,
			ForceColors:  c.Logger.Color,
		})
	case "emoji":
		log.SetFormatter(&formatter.Config{})
	case "text":
		log.SetFormatter(&log.TextFormatter{
			ForceColors: c.Logger.Color,
		})
	}

}
