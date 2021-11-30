package service

import (
	"fmt"
	"trade-bot/internal/pkg/models"

	"github.com/pkg/errors"

	"trade-bot/internal/pkg/repository"
	"trade-bot/internal/pkg/tradeAlgorithm"
	"trade-bot/internal/pkg/tradeAlgorithm/types"
	"trade-bot/internal/pkg/web"
	"trade-bot/pkg/krakenFuturesSDK"
)

var (
	ErrSendOrderServiceMethod = errors.New("send order service method")
	ErrStartTradingService    = errors.New("start trading service")
)

type KrakenOrdersManagerService struct {
	sdk    web.KrakenOrdersManager
	repo   repository.KrakenOrdersManager
	trader tradeAlgorithm.Trader
}

func NewKrakenOrdersManagerService(sdk web.KrakenOrdersManager, repo repository.KrakenOrdersManager,
	trader tradeAlgorithm.Trader) *KrakenOrdersManagerService {
	return &KrakenOrdersManagerService{sdk: sdk, repo: repo, trader: trader}
}

func (k *KrakenOrdersManagerService) SendOrder(userID int, args krakenFuturesSDK.SendOrderArguments) (models.Order, error) {
	sendStatus, err := k.sdk.SendOrder(args)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", ErrSendOrderServiceMethod, err)
	}

	order, err := k.sdk.ParseSendStatusToExecutedOrder(userID, sendStatus)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", ErrSendOrderServiceMethod, err)
	}

	if err := k.repo.CreateOrder(userID, order); err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", ErrSendOrderServiceMethod, err)
	}

	return order, nil
}

func (k *KrakenOrdersManagerService) StartTrading(userID int, details types.TradingDetails) (models.Order, error) {
	sendArgs := krakenFuturesSDK.SendOrderArguments{
		OrderType: details.OrderType,
		Symbol:    details.Symbol,
		Side:      details.Side,
		Size:      details.Size,
	}

	startOrder, err := k.SendOrder(userID, sendArgs)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", ErrStartTradingService, err)
	}

	details.BuyPrice = startOrder.Price

	if err := k.trader.StartAnalyzing(details); err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", ErrStartTradingService, err)
	}

	opositeArgs := sendArgs
	opositeArgs.ChangeToOpositeOrderSide()

	finishOrder, err := k.SendOrder(userID, opositeArgs)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", ErrStartTradingService, err)
	}

	return finishOrder, nil
}
