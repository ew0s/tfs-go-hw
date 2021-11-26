package service

import (
	"trade-bot/internal/pkg/models"
	"trade-bot/internal/pkg/repository"
	"trade-bot/internal/pkg/web"
	"trade-bot/pkg/krakenFuturesSDK"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateJWT(username string, password string) (string, error)
	GetUserIDByJWT(token string) (int, error)
	LogoutUser(token string) error
	GetUserAPIKeys(userID int) (string, string, error)
}

type KrakenOrdersManager interface {
	SendOrder(userID int, args krakenFuturesSDK.SendOrderArguments) (string, error)
}

type Service struct {
	Authorization
	KrakenOrdersManager
}

func NewService(r *repository.Repository, w *web.Web) *Service {
	return &Service{
		Authorization:       NewAuthService(r.Authorization, r.JWT),
		KrakenOrdersManager: NewKrakenOrdersManagerService(w.KrakenOrdersManager, r.KrakenOrdersManager),
	}
}
