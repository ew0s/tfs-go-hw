package service

import (
	"fmt"

	"trade-bot/internal/pkg/repository"
	"trade-bot/internal/pkg/tradeAlgorithm"
	"trade-bot/internal/pkg/tradeAlgorithm/types"
	"trade-bot/internal/pkg/web"
	"trade-bot/pkg/krakenFuturesSDK"

	"github.com/pkg/errors"
)

var (
	ErrSendOrderServiceMethod = errors.New("send order service method")
	ErrStartTradingService    = errors.New("start trading service")
)

type KrakenOrdersManagerService struct {
	sdk    web.KrakenOrdersManager
	rpeo   repository.KrakenOrdersManager
	trader tradeAlgorithm.Trader
}

func NewKrakenOrdersManagerService(sdk web.KrakenOrdersManager, rpeo repository.KrakenOrdersManager,
	trader tradeAlgorithm.Trader) *KrakenOrdersManagerService {
	return &KrakenOrdersManagerService{sdk: sdk, rpeo: rpeo, trader: trader}
}

func (k *KrakenOrdersManagerService) SendOrder(userID int, args krakenFuturesSDK.SendOrderArguments) (string, error) {
	sendStatus, err := k.sdk.SendOrder(args)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrSendOrderServiceMethod, err)
	}

	order, err := k.sdk.ParseSendStatusToExecutedOrder(userID, sendStatus)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrSendOrderServiceMethod, err)
	}

	if err := k.rpeo.CreateOrder(userID, order); err != nil {
		return "", fmt.Errorf("%s: %w", ErrSendOrderServiceMethod, err)
	}

	return order.ID, nil
}

func (k *KrakenOrdersManagerService) StartTrading(userID int, details types.TradingDetails) (string, error) {
	sendArgs := krakenFuturesSDK.SendOrderArguments{
		OrderType: details.OrderType,
		Symbol:    details.Symbol,
		Side:      details.Side,
		Size:      details.Size,
	}

	orderID, err := k.SendOrder(userID, sendArgs)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrStartTradingService, err)
	}

	order, err := k.rpeo.GetOrder(orderID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrStartTradingService, err)
	}

	details.BuyPrice = order.LimitPrice

	if err := k.trader.StartAnalyzing(details); err != nil {
		return "", fmt.Errorf("%s: %w", ErrStartTradingService, err)
	}

	opositeArgs := sendArgs
	opositeArgs.ChangeToOpositeOrderSide()

	orderID, err = k.SendOrder(userID, opositeArgs)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrStartTradingService, err)
	}

	return orderID, nil
}
