package mock

import (
	"context"

	"indexPriceCalculator/exchange"
)

type mockExchange struct {
	name       string
	subscriber exchange.PriceStreamSubscriber
}

func (m *mockExchange) GetName() string {
	return m.name
}

func (m *mockExchange) SubscribePriceStream(ctx context.Context, ticker exchange.Ticker) (chan exchange.TickerPrice, chan error) {
	return m.subscriber.SubscribePriceStream(ctx, ticker)
}

func NewExchange(cfg ExchangeConfig) exchange.Exchange {
	cfg.Validate()

	return &mockExchange{
		name:       cfg.Name,
		subscriber: NewSubscriber(cfg.TickDuration, cfg.TickersInitialValue, cfg.TickersInitialVolatility),
	}
}
