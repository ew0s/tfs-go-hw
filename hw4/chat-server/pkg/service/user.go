package service

import (
	"chat/entities"
	"chat/pkg/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SendMessageToUserByID(userID int, message entities.Message) error {
	return s.repo.SendMessageToUserByID(userID, message)
}

func (s *UserService) GetUserMessages(userID int) ([]entities.Message, error) {
	return s.repo.GetUserMessages(userID)
}
