package log

import (
	"context"

	"github.com/helmwave/helmwave/pkg/helper"
	formatter "github.com/helmwave/logrus-emoji-formatter"
	"github.com/mgutz/ansi"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var Default = &Settings{}

// Settings stores configuration for logger.
type Settings struct {
	level      string
	format     string
	color      bool
	timestamps bool
}

// Flags returns CLI flags for logger settings.
func (l *Settings) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-format",
			Usage:       "You can set: [ text | json | pad | emoji ]",
			Value:       "emoji",
			Category:    "LOGGER",
			EnvVars:     []string{"HELMWAVE_LOG_FORMAT"},
			Destination: &l.format,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "You can set: [ debug | info | warn  | fatal | panic | trace ]",
			Value:       "info",
			Category:    "LOGGER",
			EnvVars:     []string{"HELMWAVE_LOG_LEVEL", "HELMWAVE_LOG_LVL"},
			Destination: &l.level,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "on/off color",
			Value:       true,
			Category:    "LOGGER",
			EnvVars:     []string{"HELMWAVE_LOG_COLOR"},
			Destination: &l.color,
		},
		&cli.BoolFlag{
			Name:        "log-timestamps",
			Usage:       "Add timestamps to log messages",
			Value:       false,
			Category:    "LOGGER",
			EnvVars:     []string{"HELMWAVE_LOG_TIMESTAMPS"},
			Destination: &l.timestamps,
		},
	}
}

// Run initializes logger.
func (l *Settings) Run(_ *cli.Context) error {
	return l.Init()
}

// Init initializes logger and sets up hacks for other loggers (used by 3rd party libraries).
func (l *Settings) Init() error {
	// Skip various low-level k8s client errors
	// There are a lot of context deadline errors being logged
	utilruntime.ErrorHandlers = []utilruntime.ErrorHandler{ //nolint:reassign
		logKubernetesClientError,
	}

	l.setFormat()

	return l.setLevel()
}

func (l *Settings) setLevel() error {
	level, err := log.ParseLevel(l.level)
	if err != nil {
		return NewInvalidLogLevelError(l.level, err)
	}
	log.SetLevel(level)
	if level >= log.DebugLevel {
		log.SetReportCaller(true)
		helper.Helm.Debug = true
	}

	return nil
}

func (l *Settings) setFormat() {
	// Helm diff also use it
	ansi.DisableColors(!l.color)

	switch l.format {
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			PrettyPrint: true,
		})
	case "pad":
		log.SetFormatter(&log.TextFormatter{
			PadLevelText:     true,
			ForceColors:      l.color,
			FullTimestamp:    l.timestamps,
			DisableTimestamp: !l.timestamps,
		})
	case "emoji":
		cfg := &formatter.Config{
			Color: l.color,
		}

		switch {
		case !l.color && l.timestamps:
			cfg.LogFormat = "[%time%] [%lvl%]: %msg%"
		case !l.color:
			cfg.LogFormat = "[%lvl%]: %msg%"
		case l.timestamps:
			cfg.LogFormat = "[%time%] [%emoji% aka %lvl%]: %msg%"
		}

		log.SetFormatter(cfg)
	case "text":
		log.SetFormatter(&log.TextFormatter{
			ForceColors:      l.color,
			FullTimestamp:    l.timestamps,
			DisableTimestamp: !l.timestamps,
		})
	}
}

func logKubernetesClientError(ctx context.Context, err error, msg string, keysAndValues ...interface{}) {
	log.WithError(err).Trace("kubernetes client error, ", msg)
}

func (l *Settings) Format() string {
	return l.format
}
