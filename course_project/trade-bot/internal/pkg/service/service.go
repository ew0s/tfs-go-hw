package service

import (
	"trade-bot/internal/pkg/models"
	"trade-bot/internal/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateJWT(username string, password string) (string, error)
	GetUserIDByJWT(token string) (int, error)
	LogoutUser(token string) error
	GetUserAPIKeys(userID int) (string, string, error)
}

type Service struct {
	Authorization
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(r.Authorization, r.JWT),
	}
}
