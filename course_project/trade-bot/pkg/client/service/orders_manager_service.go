package service

import (
	"context"
	"fmt"
	"net/http"

	"trade-bot/pkg/client/app"
	"trade-bot/pkg/client/models"

	"github.com/pkg/errors"
)

var (
	ErrSendOrder    = errors.New("send order")
	ErrStratTrading = errors.New("strat trading")
)

type OrdersManagerService struct {
	client app.ClientActions
}

func NewOrdersManagerService(client app.ClientActions) *OrdersManagerService {
	return &OrdersManagerService{client: client}
}

func (s *OrdersManagerService) SendOrder(input models.SendOrderInput) (models.SendOrderResponse, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/orderManager/send-order", input.JWTToken, input)
	if err != nil {
		return models.SendOrderResponse{}, fmt.Errorf("%s: %w", ErrSendOrder, err)
	}

	var output models.SendOrderResponse

	resp, err := s.client.Do(req, &output)
	if err != nil {
		return models.SendOrderResponse{}, fmt.Errorf("%s: %w", ErrSendOrder, err)
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 400) {
		return models.SendOrderResponse{}, fmt.Errorf("%s: %s: %s", ErrSendOrder, resp.Status, output.Message)
	}

	return output, err
}

func (s *OrdersManagerService) StartTrading(ctx context.Context, input models.StartTradingInput) (<-chan *models.StartTradingResponse, <-chan error, error) {
	req, err := s.client.NewWsRequest("/orderManager/ws/start-trade", input.JWTToken)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", ErrStratTrading, err)
	}

	conn, err := s.client.DoWS(req, input)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", ErrStratTrading, err)
	}

	var output models.StartTradingResponse

	respCh, errCh := s.client.LoopOverWS(ctx, conn, &output)

	tradingRespCh := make(chan *models.StartTradingResponse)
	go func() {
		defer close(tradingRespCh)

		for val := range respCh {
			tradingRespCh <- val.(*models.StartTradingResponse)
		}
	}()

	return tradingRespCh, errCh, nil
}
