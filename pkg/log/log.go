package log

import (
	"fmt"

	"github.com/bombsimon/logrusr/v2"
	"github.com/helmwave/helmwave/pkg/helper"
	formatter "github.com/helmwave/logrus-emoji-formatter"
	"github.com/mgutz/ansi"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/werf/logboek"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
)

// Settings stores configuration for logger.
type Settings struct {
	level      string
	format     string
	color      bool
	timestamps bool
	width      int
}

// Flags returns CLI flags for logger settings.
func (l *Settings) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-format",
			Usage:       "You can set: [ text | json | pad | emoji ]",
			Value:       "emoji",
			EnvVars:     []string{"HELMWAVE_LOG_FORMAT"},
			Destination: &l.format,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "You can set: [ debug | info | warn  | fatal | panic | trace ]",
			Value:       "info",
			EnvVars:     []string{"HELMWAVE_LOG_LEVEL", "HELMWAVE_LOG_LVL"},
			Destination: &l.level,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Force color",
			Value:       true,
			EnvVars:     []string{"HELMWAVE_LOG_COLOR"},
			Destination: &l.color,
		},
		&cli.IntFlag{
			Name:        "kubedog-log-width",
			Usage:       "Set kubedog max log line width",
			Value:       140,
			EnvVars:     []string{"HELMWAVE_KUBEDOG_LOG_WIDTH"},
			Destination: &l.width,
		},
		&cli.BoolFlag{
			Name:        "log-timestamps",
			Usage:       "Add timestamps to log messages",
			Value:       false,
			EnvVars:     []string{"HELMWAVE_LOG_TIMESTAMPS"},
			Destination: &l.timestamps,
		},
	}
}

// Run initializes logger.
func (l *Settings) Run(c *cli.Context) error {
	return l.Init()
}

// Init initializes logger and sets up hacks for other loggers (used by 3rd party libraries).
func (l *Settings) Init() error {
	// Skip various low-level k8s client errors
	// There are a lot of context deadline errors being logged
	utilruntime.ErrorHandlers = []func(error){
		logKubernetesClientError,
	}
	klog.SetLogger(logrusr.New(log.StandardLogger()))

	if l.width > 0 {
		logboek.DefaultLogger().Streams().SetWidth(l.width)
	}

	l.setFormat()

	return l.setLevel()
}

func (l *Settings) setLevel() error {
	level, err := log.ParseLevel(l.level)
	if err != nil {
		return fmt.Errorf("failed to parse log level %s: %w", l.level, err)
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

		if !l.color && l.timestamps { // nolint:gocritic
			cfg.LogFormat = "[%time%] [%lvl%]: %msg%"
		} else if !l.color {
			cfg.LogFormat = "[%lvl%]: %msg%"
		} else if l.timestamps {
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

func logKubernetesClientError(err error) {
	log.WithError(err).Debug("kubernetes client error")
}
