package kubedog

import (
	"time"
)

type Config struct {
	StatusInterval time.Duration
	Timeout        time.Duration
	StartDelay     time.Duration
}
