package index

import (
	"context"
	"fmt"
	"os"

	"indexPriceCalculator/exchange"
)

type exchangeEntity struct {
	parent *Subscriber
	api    exchange.Exchange
	// Latest price for particular ticker
	lastPrices map[exchange.Ticker]*exchange.TickerPrice
}

// Start listening for particular ticker in this exchange
func (e *exchangeEntity) start(ctx context.Context, ticker exchange.Ticker) {
	data, subscribedErrs := e.api.SubscribePriceStream(ctx, ticker)

	errs := e.parent.errors[ticker]

	name := e.api.GetName()

	for {
		select {
		// handle context shutdown
		case _ = <-ctx.Done():
			return
		// update ticker's last price for current exchange
		case price := <-data:
			e.lastPrices[ticker] = &price
			if os.Getenv("DEBUG") == "1" {
				fmt.Printf("|%-10v|	Timestamp: %d\tPrice: %f\tVolume:\t%f\n", name, price.Time.Unix(), price.Price, price.Volume)
			}
			break
		// send error to index worker
		case err := <-subscribedErrs:
			errs <- fmt.Errorf("%s | %v", name, err)
			if e.parent.exitOnError {
				close(data)
				close(subscribedErrs)
				e.lastPrices[ticker] = nil
				return

			}

		}

	}
}
