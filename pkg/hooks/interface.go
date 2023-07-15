package hooks

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type Hook interface {
	Run(context.Context) error
	Log() *log.Entry
}
