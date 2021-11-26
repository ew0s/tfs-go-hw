package service

import (
	"fmt"
	"net/http"
	"trade-bot/pkg/client/app"
	"trade-bot/pkg/client/models"

	"github.com/pkg/errors"
)

var (
	ErrSendOrder = errors.New("send order")
)

type OrdersManagerService struct {
	client app.ClientActions
}

func NewOrdersManagerService(client app.ClientActions) *OrdersManagerService {
	return &OrdersManagerService{client: client}
}

func (s *OrdersManagerService) SendOrder(input models.SendOrderInput) (models.SendOrderResponse, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/orderManager/send-order", input.JwtToken, input)
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
