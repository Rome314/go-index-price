package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"indexPriceCalculator/exchange"
	"indexPriceCalculator/index"
	"indexPriceCalculator/mock"
)

func main() {

	// If 1 then exchanges live prices are shown
	os.Setenv("DEBUG", "0")
	indexLogPrefix := ""
	if os.Getenv("DEBUG") == "1" {
		indexLogPrefix = "\n"
	}

	exchangesTickerDur := time.Second * 3

	bitmex := mock.NewExchange(mock.ExchangeConfig{Name: "Bitmex", TickDuration: exchangesTickerDur})
	binance := mock.NewExchange(mock.ExchangeConfig{Name: "Binance", TickDuration: exchangesTickerDur})
	ftx := mock.NewExchange(mock.ExchangeConfig{Name: "FTX", TickDuration: exchangesTickerDur})

	compositeIndex := index.NewSubscriber(index.Config{ExitOnError: false, TickerDuration: time.Second * 60},
		bitmex, binance, ftx)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	data, errs := compositeIndex.SubscribePriceStream(ctx, exchange.BTCUSDTicker)

	for {
		select {
		case ticker := <-data:
			fmt.Printf(
				"%s|%-10v|\tTimestamp: %d\tPrice: %f\tVolume:\t%f\n\n",
				indexLogPrefix,
				"INDEX",
				ticker.Time.Unix(),
				ticker.Price,
				ticker.Volume,
			)
			break
		case err := <-errs:
			fmt.Printf("Error: %s\n", err.Error())

		}
	}
}
