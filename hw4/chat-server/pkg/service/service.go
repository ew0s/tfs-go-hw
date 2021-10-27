package service

import (
	"chat/entities"
	"chat/pkg/repository"
)

type Authorization interface {
	CreateUser(user entities.User) (int, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Message interface {
	Create(usrID int, message entities.Message) error
	GetGlobalMessages() []entities.Message
}

type User interface {
	SendMessageToUserByID(userID int, message entities.Message) error
	GetUserMessages(userID int) ([]entities.Message, error)
}

type Service struct {
	Authorization
	Message
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Message:       NewMessageService(repos.Message),
		User:          NewUserService(repos.User),
	}
}
