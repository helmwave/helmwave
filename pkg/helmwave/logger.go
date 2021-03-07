package helmwave

import (
	"github.com/bombsimon/logrusr"
	log "github.com/sirupsen/logrus"
	"github.com/zhilyaev/helmwave/pkg/formatter"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
)

type Log struct {
	//Engine *log.Logger
	Level  string
	Format string
	Color  bool
}

func (c *Config) InitLogger() error {
	// Skip various low-level k8s client errors
	// There are a lot of context deadline errors being logged
	utilruntime.ErrorHandlers = []func(error){
		logKubernetesClientError,
	}
	klog.SetLogger(logrusr.NewLogger(log.StandardLogger()))

	c.InitLoggerFormat()
	return c.InitLoggerLevel()
}

func (c *Config) InitLoggerLevel() error {
	level, err := log.ParseLevel(c.Logger.Level)
	if err != nil {
		return err
	}
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
		log.SetFormatter(&formatter.Config{
			Color: c.Logger.Color,
		})
	case "text":
		log.SetFormatter(&log.TextFormatter{
			ForceColors: c.Logger.Color,
		})
	}

}

func logKubernetesClientError(err error) {
	log.Debugf("kubernetes client error %q", err.Error())
}
