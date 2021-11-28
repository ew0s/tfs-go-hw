package tradeAlgorithm

import (
	"trade-bot/internal/pkg/tradeAlgorithm/algorithms"
	"trade-bot/internal/pkg/tradeAlgorithm/types"
	"trade-bot/internal/pkg/web"
)

type Trader interface {
	StartAnalyzing(details types.TradingDetails) error
}

type TradeAlgorithm struct {
	Trader
}

func NewTradeAlgorithm(w *web.Web) *TradeAlgorithm {
	return &TradeAlgorithm{Trader: algorithms.NewStopLossTakeProfitAlgo(w.KrakenAnalyzer)}
}
