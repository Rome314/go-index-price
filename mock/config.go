package mock

import (
	"fmt"
	"math"
	"time"
)

type ExchangeConfig struct {
	Name                     string
	TickDuration             time.Duration
	TickersInitialValue      float64
	TickersInitialVolatility float64
}

// Setting default fields if they aren't present
func (cfg *ExchangeConfig) Validate() {
	if cfg.TickersInitialValue == 0 {
		cfg.TickersInitialValue = 1000
	}
	if cfg.TickersInitialVolatility == 0 {
		cfg.TickersInitialVolatility = 10
	}
	if cfg.Name == "" {
		cfg.Name = fmt.Sprintf("Exchange | %f", math.Mod(float64(time.Now().Second()), 1000000000))
	}
	if cfg.TickDuration.Nanoseconds() == 0 {
		cfg.TickDuration = time.Second
	}

}
