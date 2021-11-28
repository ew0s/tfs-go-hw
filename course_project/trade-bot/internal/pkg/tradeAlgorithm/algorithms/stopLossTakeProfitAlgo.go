package algorithms

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"trade-bot/internal/pkg/tradeAlgorithm/types"
	"trade-bot/internal/pkg/web"
	"trade-bot/pkg/krakenFuturesWSSDK"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrStartAnalyzing     = errors.New("strat analyzing")
	ErrUnableToGetCandles = errors.New("unable to get candles")
)

type StopLossTakeProfitAlgo struct {
	krakenWebsocketSDK web.KrakenAnalyzer
}

func NewStopLossTakeProfitAlgo(krakenAnalyzer web.KrakenAnalyzer) *StopLossTakeProfitAlgo {
	return &StopLossTakeProfitAlgo{
		krakenWebsocketSDK: krakenAnalyzer,
	}
}

func (a *StopLossTakeProfitAlgo) StartAnalyzing(details types.TradingDetails) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	candles, err := a.krakenWebsocketSDK.LookForCandles(ctx, krakenFuturesWSSDK.OneMinuteCandlesFeed, []string{details.Symbol})
	if err != nil {
		return fmt.Errorf("%s: %w", ErrStartAnalyzing, err)
	}

	for candle := range candles {
		log.Info(candle)
		price, err := strconv.ParseFloat(candle.Close, 64)
		log.Infof("\nclose: %s\nto_take_profit:%f\nto_stop_loss:%f\n",
			candle.Close,
			math.Abs(details.BuyPrice+details.TakeProfitBorder-price),
			math.Abs(details.BuyPrice-details.StopLossBorder-price))

		if err != nil {
			return fmt.Errorf("%s: %w", ErrStartAnalyzing, err)
		}

		if price > details.BuyPrice+details.TakeProfitBorder {
			return nil
		}
		if price < details.BuyPrice-details.StopLossBorder {
			return nil
		}
	}

	return fmt.Errorf("%s: %s", ErrStartAnalyzing, ErrUnableToGetCandles)
}
