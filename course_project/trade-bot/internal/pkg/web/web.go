package web

import (
	"trade-bot/internal/pkg/web/webSDK"
	"trade-bot/pkg/krakenFuturesSDK"
)

type KrakenOrdersManager interface {
	SendOrder(args krakenFuturesSDK.SendOrderArguments) (krakenFuturesSDK.SendStatus, error)
	EditOrder(args krakenFuturesSDK.EditOrderArguments) (krakenFuturesSDK.EditStatus, error)
	CancelOrder(args krakenFuturesSDK.CancelOrderArguments) (krakenFuturesSDK.CancelStatus, error)
	CancelAllOrders(symbol string) (krakenFuturesSDK.CancelAllStatus, error)
}

type Web struct {
	KrakenOrdersManager
}

func NewWeb(krakenWebSDK *krakenFuturesSDK.API) *Web {
	return &Web{
		KrakenOrdersManager: webSDK.NewKrakenOrdersManagerWebSDK(krakenWebSDK),
	}
}
