package kubedog

import (
	"time"
)

// Config is config for kubedog library.
type Config struct {
	Enabled        bool
	TrackGeneric   bool
	StatusInterval time.Duration
	Timeout        time.Duration
	StartDelay     time.Duration
	LogWidth       int
}
