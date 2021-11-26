package service

import (
	"fmt"
	"trade-bot/internal/pkg/models"
	"trade-bot/internal/pkg/repository"
	"trade-bot/internal/pkg/web"
	"trade-bot/pkg/krakenFuturesSDK"

	"github.com/pkg/errors"
)

var (
	ErrSendOrderServiceMethod = errors.New("send order service method")
	ErrUnknownSendStatusType  = errors.New("unknown send status type")
)

type KrakenOrdersManagerService struct {
	sdk  web.KrakenOrdersManager
	rpeo repository.KrakenOrdersManager
}

func NewKrakenOrdersManagerService(sdk web.KrakenOrdersManager, rpeo repository.KrakenOrdersManager) *KrakenOrdersManagerService {
	return &KrakenOrdersManagerService{sdk: sdk, rpeo: rpeo}
}

func (k *KrakenOrdersManagerService) SendOrder(userID int, args krakenFuturesSDK.SendOrderArguments) (string, error) {
	sendStatus, err := k.sdk.SendOrder(args)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrSendOrderServiceMethod, err)
	}

	order, err := parseSendStatusToExecutedOrder(userID, sendStatus)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrSendOrderServiceMethod, err)
	}

	if err := k.rpeo.CreateOrder(userID, order); err != nil {
		return "", fmt.Errorf("%s: %w", ErrSendOrderServiceMethod, err)
	}

	return order.ID, nil
}

func parseSendStatusToExecutedOrder(userID int, sendStatus krakenFuturesSDK.SendStatus) (models.Order, error) {
	orderEvent := sendStatus.OrderEvents[0]

	if orderEvent.Type == "EXECUTION" {
		return models.Order{
			ID:                  orderEvent.OrderPriorExecution.OrderID,
			UserID:              userID,
			ClientOrderID:       orderEvent.OrderPriorExecution.CliOrderID,
			Type:                orderEvent.Type,
			Symbol:              orderEvent.OrderPriorExecution.Symbol,
			Quantity:            orderEvent.OrderPriorExecution.Quantity,
			Side:                orderEvent.OrderPriorExecution.Side,
			LimitPrice:          orderEvent.OrderPriorExecution.LimitPrice,
			Filled:              orderEvent.OrderPriorExecution.Filled,
			Timestamp:           orderEvent.OrderPriorExecution.Timestamp,
			LastUpdateTimestamp: orderEvent.OrderPriorExecution.LastUpdateTimestamp,
		}, nil
	}

	return models.Order{}, ErrUnknownSendStatusType
}
