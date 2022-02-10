package mock

import (
	"context"
	"math/rand"
	"time"

	"indexPriceCalculator/exchange"
)

type mockStream struct {
	ticker    exchange.Ticker
	startTime time.Time
	// Initial price for coin, configurable
	initialValue float64
	// Current price for coin
	currentValue float64
	// Amount by which current price changes randomly up or down, configurable
	volatility float64
	// At every tick new price sending to the channel
	timeTicker *time.Ticker
}

func NewSubscriber(tickDur time.Duration, initial, volatility float64) exchange.PriceStreamSubscriber {
	return &mockStream{
		initialValue: initial,
		currentValue: initial,
		volatility:   volatility,
		timeTicker:   time.NewTicker(tickDur),
	}
}

func (m *mockStream) SubscribePriceStream(ctx context.Context, ticker exchange.Ticker) (chan exchange.TickerPrice, chan error) {
	m.ticker = ticker
	out := make(chan exchange.TickerPrice)
	err := make(chan error)

	go m.start(ctx, out, err)
	return out, err

}
func (m *mockStream) start(ctx context.Context, out chan exchange.TickerPrice, err chan error) {
	m.startTime = time.Now()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		select {
		case _ = <-ctx.Done():
			close(err)
			close(out)
			return
		case timestamp := <-m.timeTicker.C:
			multiplier := -1.0
			if r.Intn(2) == 1 {
				multiplier = 1.0
			}

			m.currentValue += multiplier * m.volatility

			out <- exchange.TickerPrice{
				Ticker: m.ticker,
				Time:   timestamp,
				Price:  m.currentValue,
				// Volume always random
				Volume: r.Float64() * 1000000,
			}

		}
	}

}
