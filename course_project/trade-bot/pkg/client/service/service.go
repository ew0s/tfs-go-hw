package service

import (
	"trade-bot/pkg/client/app"
	"trade-bot/pkg/client/models"
)

type Authorization interface {
	SignUp(input models.SignUpInput) (models.SignUpResponse, error)
	SignIn(input models.SignInInput) (models.SignInResponse, error)
	Logout(input models.LogoutInput) (models.LogoutResponse, error)
}

type OrdersManagee interface {
	SendOrder(input models.SendOrderInput) (models.SendOrderResponse, error)
}

type Service struct {
	Authorization
	OrdersManagee
}

func NewService(client app.ClientActions) *Service {
	return &Service{
		Authorization: NewAuthService(client),
		OrdersManagee: NewOrdersManagerService(client),
	}
}
