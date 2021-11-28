package service

import (
	"context"

	"trade-bot/pkg/client/app"
	"trade-bot/pkg/client/models"
)

type Authorization interface {
	SignUp(input models.SignUpInput) (models.SignUpResponse, error)
	SignIn(input models.SignInInput) (models.SignInResponse, error)
	Logout(input models.LogoutInput) (models.LogoutResponse, error)
}

type OrdersManager interface {
	SendOrder(input models.SendOrderInput) (models.SendOrderResponse, error)
	StartTrading(ctx context.Context, input models.StartTradingInput) (<-chan *models.StartTradingResponse, <-chan error, error)
}

type Service struct {
	Authorization
	OrdersManager
}

func NewService(client app.ClientActions) *Service {
	return &Service{
		Authorization: NewAuthService(client),
		OrdersManager: NewOrdersManagerService(client),
	}
}
