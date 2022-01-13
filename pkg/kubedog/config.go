package kubedog

import (
	"time"
)

// Config is config for kubedog library.
type Config struct {
	StatusInterval time.Duration
	Timeout        time.Duration
	StartDelay     time.Duration
}
