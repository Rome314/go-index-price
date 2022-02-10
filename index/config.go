package index

import (
	"time"
)

type Config struct {
	// If this field disabled then subscribers just send errors to channel, otherwise - stopping
	ExitOnError bool
	// Period when index calculated
	TickerDuration time.Duration
}

// Setting default values for missed important fields
func (cfg *Config) Validate() {
	if cfg.TickerDuration.Nanoseconds() == 0 {
		cfg.TickerDuration = time.Minute
	}
}
