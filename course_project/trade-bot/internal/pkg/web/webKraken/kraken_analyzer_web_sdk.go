package webKraken

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"trade-bot/pkg/krakenFuturesWSSDK"
)

var (
	ErrConvertTradeDataToCandle = errors.New("convert trade data to candle")
	ErrLookForCandles           = errors.New("look for candles")
)

type KrakenAnalyzerWebSDK struct {
	krakenWebsocketAPI *krakenFuturesWSSDK.WSAPI
}

func NewKrakenAnalyzerWebSDK(krakenWebsocketAPI *krakenFuturesWSSDK.WSAPI) *KrakenAnalyzerWebSDK {
	return &KrakenAnalyzerWebSDK{krakenWebsocketAPI: krakenWebsocketAPI}
}

func (k *KrakenAnalyzerWebSDK) LookForCandles(ctx context.Context, feed string, productsIDs []string) (<-chan krakenFuturesWSSDK.Candle, error) {
	tradeDataCh, err := k.krakenWebsocketAPI.CandlesTrade(ctx, feed, productsIDs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrLookForCandles, err)
	}

	candleCh, errCh := convertTradeDataToCandle(tradeDataCh)
	go logErrors(errCh)
	return filterCandles(candleCh), nil
}

func logErrors(errs <-chan error) {
	for err := range errs {
		log.Warn(err)
	}
}

func convertTradeDataToCandle(tradeData <-chan *krakenFuturesWSSDK.CandlesTradeData) (<-chan krakenFuturesWSSDK.Candle, <-chan error) {
	errCh := make(chan error, 1)
	candlesChan := make(chan krakenFuturesWSSDK.Candle)

	go func() {
		defer close(candlesChan)
		defer close(errCh)

		for data := range tradeData {
			if data.Feed == "error" {
				errCh <- fmt.Errorf("%s: error feed sended", ErrConvertTradeDataToCandle)
				continue
			}

			candlesChan <- data.Candle
		}
	}()

	return candlesChan, errCh
}

func filterCandles(candles <-chan krakenFuturesWSSDK.Candle) <-chan krakenFuturesWSSDK.Candle {
	candlesChan := make(chan krakenFuturesWSSDK.Candle)

	go func() {
		defer close(candlesChan)

		var lastUpdateTime *int

		for candle := range candles {
			if lastUpdateTime == nil {
				t := candle.Time
				lastUpdateTime = &t
				candlesChan <- candle
				continue
			}
			if *lastUpdateTime == candle.Time {
				continue
			}

			*lastUpdateTime = candle.Time
			candlesChan <- candle
		}
	}()

	return candlesChan
}
