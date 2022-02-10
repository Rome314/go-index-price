package exchange

import (
	"context"
	"time"
)

type Exchange interface {
	GetName() string
	PriceStreamSubscriber
}

type Ticker string

const (
	BTCUSDTicker Ticker = "BTC_USD"
)

type TickerPrice struct {
	Ticker Ticker
	Time   time.Time
	// Changed from string to float for comfortable calculations
	Price float64
	// Added for index algorithm, BTW in common exchanges  provide this information anyway
	Volume float64 // The trading volume of the symbol in the given time period.
}

type PriceStreamSubscriber interface {
	// Added context to parameters for graceful stopping streams
	SubscribePriceStream(context.Context, Ticker) (chan TickerPrice, chan error)
}
