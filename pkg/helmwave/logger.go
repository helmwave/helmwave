package helmwave

import log "github.com/sirupsen/logrus"

type Log struct {
	//Engine *log.Logger
	Level  string
	Format string
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
			ForceColors:  true,
		})
	case "text":
		log.SetFormatter(&log.TextFormatter{
			ForceColors: true,
		})
	}
}
