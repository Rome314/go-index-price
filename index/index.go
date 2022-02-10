package index

import (
	"context"
	"fmt"
	"time"

	"indexPriceCalculator/exchange"
)

type Subscriber struct {
	exchanges   map[string]*exchangeEntity
	exitOnError bool
	timeTicker  *time.Ticker
	errors      map[exchange.Ticker]chan error
}

func NewSubscriber(cfg Config, exchanges ...exchange.Exchange) exchange.PriceStreamSubscriber {
	cfg.Validate()

	s := &Subscriber{
		exitOnError: cfg.ExitOnError,
		exchanges:   map[string]*exchangeEntity{},
		timeTicker:  time.NewTicker(cfg.TickerDuration),
		errors:      map[exchange.Ticker]chan error{},
	}

	for _, e := range exchanges {
		s.exchanges[e.GetName()] = &exchangeEntity{
			parent:     s,
			api:        e,
			lastPrices: map[exchange.Ticker]*exchange.TickerPrice{},
		}
	}
	return s
}

func (s *Subscriber) SubscribePriceStream(ctx context.Context, ticker exchange.Ticker) (chan exchange.TickerPrice, chan error) {
	s.errors[ticker] = make(chan error)

	streamContex, streamsCancel := context.WithCancel(ctx)

	for _, e := range s.exchanges {
		go e.start(streamContex, ticker)
	}

	out := make(chan exchange.TickerPrice)
	errs := make(chan error)

	go s.start(ctx, streamsCancel, ticker, out, errs)
	return out, errs
}

func (s *Subscriber) start(ctx context.Context, cancel context.CancelFunc, ticker exchange.Ticker, out chan exchange.TickerPrice, errs chan error) {
	defer cancel()
	for {
		select {
		case _ = <-ctx.Done():
			close(out)
			return
		case timestamp := <-s.timeTicker.C:

			totalVolume := 0.0
			pricesByExchange := map[string]*exchange.TickerPrice{}

			for name, exchng := range s.exchanges {
				lastPrice, ok := exchng.lastPrices[ticker]
				if !ok || lastPrice == nil {
					continue
				}

				totalVolume += lastPrice.Volume
				pricesByExchange[name] = lastPrice
			}

			sumWeights := 0.0
			sumPrices := 0.0
			for _, price := range pricesByExchange {
				weight := price.Volume / totalVolume

				sumPrices += price.Price * weight
				sumWeights += weight
			}

			out <- exchange.TickerPrice{
				Ticker: ticker,
				Time:   timestamp,
				Price:  sumPrices / sumWeights,
				Volume: totalVolume,
			}
		case err := <-s.errors[ticker]:
			errs <- err
			if s.exitOnError {
				errs <- fmt.Errorf("exitOnerror is on, closing")
				close(out)
				close(errs)
				return
			}

		}

	}
}
