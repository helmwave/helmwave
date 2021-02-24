package kubedog

import (
	"time"
)

type Config struct {
	Enabled        bool
	StatusInterval time.Duration
	Timeout        time.Duration
	StartDelay     time.Duration
}
