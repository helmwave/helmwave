package tracker

import (
	"time"
)

// Config is the configuration for live resource tracking.
type Config struct {
	Enabled        bool
	TrackGeneric   bool
	Logs           bool
	StatusInterval time.Duration
	Timeout        time.Duration
	StartDelay     time.Duration
	LogWidth       int // Maximum streamed log line width, 0 means unlimited.
}
