package log

import (
	log "github.com/sirupsen/logrus"
)

// LoggerGetter is an interface that provides logger.
type LoggerGetter interface {
	Logger() *log.Entry
}
