package kubedog

import (
	"time"
)

// Config is config for kubedog library.
type Config struct {
	Enabled        bool
	StatusInterval time.Duration
	Timeout        time.Duration
	StartDelay     time.Duration
	LogWidth       int
}
